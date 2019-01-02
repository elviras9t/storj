// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"storj.io/storj/pkg/accounting/rollup"
	"storj.io/storj/pkg/accounting/tally"
	"storj.io/storj/pkg/audit"
	"storj.io/storj/pkg/bwagreement"
	"storj.io/storj/pkg/cfgstruct"
	"storj.io/storj/pkg/datarepair/checker"
	"storj.io/storj/pkg/datarepair/repairer"
	"storj.io/storj/pkg/discovery"
	"storj.io/storj/pkg/kademlia"
	"storj.io/storj/pkg/miniogw"
	"storj.io/storj/pkg/overlay"
	"storj.io/storj/pkg/piecestore/psserver"
	"storj.io/storj/pkg/pointerdb"
	"storj.io/storj/pkg/process"
	"storj.io/storj/pkg/satellite/satelliteweb"
	"storj.io/storj/pkg/server"
	"storj.io/storj/pkg/statdb"
	"storj.io/storj/pkg/utils"
	"storj.io/storj/satellite/satellitedb"
)

const (
	storagenodeCount = 10
)

// Satellite is for configuring client
type Satellite struct {
	Server      server.Config
	Kademlia    kademlia.SatelliteConfig
	PointerDB   pointerdb.Config
	Overlay     overlay.Config
	Checker     checker.Config
	Repairer    repairer.Config
	Audit       audit.Config
	BwAgreement bwagreement.Config
	Web         satelliteweb.Config
	Database    string `help:"satellite database connection string" default:"sqlite3://$CONFDIR/master.db"`
	Discovery   discovery.Config
	Tally       tally.Config
	Rollup      rollup.Config
	StatDB      statdb.Config
}

// StorageNode is for configuring storage nodes
type StorageNode struct {
	Server   server.Config
	Kademlia kademlia.StorageNodeConfig
	Storage  psserver.Config
}

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run all servers",
		RunE:  cmdRun,
	}

	runCfg struct {
		Satellite    Satellite
		StorageNodes [storagenodeCount]StorageNode
		Uplink       miniogw.Config
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	cfgstruct.Bind(runCmd.Flags(), &runCfg, cfgstruct.ConfDir(defaultConfDir))
}

func cmdRun(cmd *cobra.Command, args []string) (err error) {
	ctx := process.Ctx(cmd)
	defer mon.Task()(&ctx)(&err)

	errch := make(chan error, len(runCfg.StorageNodes)+2)
	// start mini redis
	m := miniredis.NewMiniRedis()
	m.RequireAuth("abc123")

	if err = m.StartAddr(":6378"); err != nil {
		errch <- err
	} else {
		defer m.Close()
	}

	// start satellite
	go func() {
		_, _ = fmt.Printf("Starting satellite on %s\n",
			runCfg.Satellite.Server.PublicAddress)

		if runCfg.Satellite.Audit.SatelliteAddr == "" {
			runCfg.Satellite.Audit.SatelliteAddr = runCfg.Satellite.Server.PublicAddress
		}

		if runCfg.Satellite.Web.SatelliteAddr == "" {
			runCfg.Satellite.Web.SatelliteAddr = runCfg.Satellite.Server.PublicAddress
		}

		database, err := satellitedb.New(runCfg.Satellite.Database)
		if err != nil {
			errch <- errs.New("Error starting master database on satellite: %+v", err)
			return
		}

		err = database.CreateTables()
		if err != nil {
			errch <- errs.New("Error creating tables for master database on satellite: %+v", err)
			return
		}

		//nolint ignoring context rules to not create cyclic dependency, will be removed later
		ctx = context.WithValue(ctx, "masterdb", database)

		// Run satellite
		errch <- runCfg.Satellite.Server.Run(ctx,
			runCfg.Satellite.Kademlia,
			runCfg.Satellite.Audit,
			runCfg.Satellite.Overlay,
			runCfg.Satellite.Discovery,
			runCfg.Satellite.PointerDB,
			runCfg.Satellite.Checker,
			runCfg.Satellite.Repairer,
			runCfg.Satellite.BwAgreement,
			runCfg.Satellite.Web,
			runCfg.Satellite.Tally,
			runCfg.Satellite.Rollup,
			runCfg.Satellite.StatDB,
		)
	}()

	// hack-fix t oensure that satellite gets up and running before starting storage nodes
	time.Sleep(2 * time.Second)

	// start the storagenodes
	for i, v := range runCfg.StorageNodes {
		go func(i int, v StorageNode) {
			identity, err := v.Server.Identity.Load()
			if err != nil {
				return
			}

			address := v.Server.PublicAddress
			storagenode := fmt.Sprintf("%s:%s", identity.ID.String(), address)

			_, _ = fmt.Printf("Starting storage node %d %s (kad on %s)\n", i, storagenode, address)
			errch <- v.Server.Run(ctx, v.Kademlia, v.Storage)
		}(i, v)
	}
	// start s3 uplink
	go func() {
		_, _ = fmt.Printf("Starting s3-gateway on %s\nAccess key: %s\nSecret key: %s\n",
			runCfg.Uplink.Server.Address,
			runCfg.Uplink.Minio.AccessKey,
			runCfg.Uplink.Minio.SecretKey)
		errch <- runCfg.Uplink.Run(ctx)
	}()

	return utils.CollectErrors(errch, 5*time.Second)
}
