// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package overlay_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"storj.io/storj/internal/testcontext"
	"storj.io/storj/internal/testplanet"
	"storj.io/storj/pkg/overlay"
	"storj.io/storj/pkg/pb"
)

func TestServer(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	planet, err := testplanet.New(t, 1, 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Check(planet.Shutdown)

	planet.Start(ctx)
	// we wait a second for all the nodes to complete bootstrapping off the satellite
	time.Sleep(2 * time.Second)

	satellite := planet.Satellites[0]
	server := satellite.Overlay.Endpoint
	// TODO: handle cleanup

	{ // FindStorageNodes
		result, err := server.FindStorageNodes(ctx, &pb.FindStorageNodesRequest{
			Opts: &pb.OverlayOptions{Amount: 2},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Nodes, 2)
	}

	{ // Lookup
		result, err := server.Lookup(ctx, &pb.LookupRequest{
			NodeId: planet.StorageNodes[0].ID(),
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, result.Node.Address.Address, planet.StorageNodes[0].Addr())
	}

	{ // BulkLookup
		result, err := server.BulkLookup(ctx, &pb.LookupRequests{
			LookupRequest: []*pb.LookupRequest{
				{NodeId: planet.StorageNodes[0].ID()},
				{NodeId: planet.StorageNodes[1].ID()},
				{NodeId: planet.StorageNodes[2].ID()},
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.LookupResponse, 3)

		for i, resp := range result.LookupResponse {
			if assert.NotNil(t, resp.Node) {
				assert.Equal(t, resp.Node.Address.Address, planet.StorageNodes[i].Addr())
			}
		}
	}
}

func TestNewNodeFiltering(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	var totalNodes int
	totalNodes = 10

	planet, err := testplanet.New(t, 1, totalNodes, 1)
	if err != nil {
		t.Fatal(err)
	}

	planet.Start(ctx)

	defer ctx.Check(planet.Shutdown)

	// we wait a second for all the nodes to complete bootstrapping off the satellite
	time.Sleep(5 * time.Second)

	satellite := planet.Satellites[0]

	// This sets a reputable audit count for a certain number of nodes.
	for i, node := range planet.StorageNodes {
		for j := 0; j < i; j++ {
			_, err := satellite.DB.StatDB().UpdateAuditSuccess(ctx, node.ID(), true)
			assert.NoError(t, err)
		}
	}

	for _, tt := range []struct {
		name                  string
		newNodeAuditThreshold int64
		newNodePercentage     float64
		requestedNodeAmt      int64
		expectedResultLength  int
	}{
		{
			name:                  "case: all reputable nodes, only reputable nodes requested",
			newNodeAuditThreshold: 0,
			newNodePercentage:     0,
			requestedNodeAmt:      5,
			expectedResultLength:  5,
		},
		// {
		// 	name:                  "case: all reputable nodes, reputable and new nodes requested",
		// 	newNodeAuditThreshold: 0,
		// 	newNodePercentage:     1,
		// 	requestedNodeAmt:      5,
		// 	expectedResultLength:  5,
		// },
		// {
		// 	name:                  "case: all reputable nodes except one, reputable and new nodes requested",
		// 	newNodeAuditThreshold: 1,
		// 	newNodePercentage:     1,
		// 	requestedNodeAmt:      5,
		// 	expectedResultLength:  6,
		// },
		// {
		// 	name:                  "case: 50-50 reputable and new nodes, reputable and new nodes requested (new node % 1)",
		// 	newNodeAuditThreshold: 5,
		// 	newNodePercentage:     1,
		// 	requestedNodeAmt:      2,
		// 	expectedResultLength:  4,
		// },
		// {
		// 	name:                  "case: 50-50 reputable and new nodes, reputable and new nodes requested (new node % .5)",
		// 	newNodeAuditThreshold: 5,
		// 	newNodePercentage:     0.5,
		// 	requestedNodeAmt:      4,
		// 	expectedResultLength:  6,
		// },
		// {
		// 	name:                  "case: all new nodes except one, reputable and new nodes requested (happy path)",
		// 	newNodeAuditThreshold: 9,
		// 	newNodePercentage:     1,
		// 	requestedNodeAmt:      1,
		// 	expectedResultLength:  2,
		// },
		// {
		// 	name:                  "case: all new nodes except one, reputable and new nodes requested (not happy path)",
		// 	newNodeAuditThreshold: 9,
		// 	newNodePercentage:     1,
		// 	requestedNodeAmt:      2,
		// 	expectedResultLength:  3,
		// },
		// {
		// 	name:                  "case: all new nodes, reputable and new nodes requested",
		// 	newNodeAuditThreshold: 50,
		// 	newNodePercentage:     1,
		// 	requestedNodeAmt:      2,
		// 	expectedResultLength:  2,
		// },
	} {

		nodeSelectionConfig := &overlay.NodeSelectionConfig{
			UptimeCount:           0,
			UptimeRatio:           0,
			AuditSuccessRatio:     0,
			AuditCount:            0,
			NewNodeAuditThreshold: tt.newNodeAuditThreshold,
			NewNodePercentage:     tt.newNodePercentage,
		}

		server := overlay.NewServer(satellite.Log.Named("overlay"), satellite.Overlay.Service, nodeSelectionConfig)

		for i := range planet.StorageNodes {

			address := "127.0.0.1:555" + strconv.Itoa(i)

			n := &pb.Node{
				Id:      planet.StorageNodes[i].ID(),
				Address: &pb.NodeAddress{Address: address},
			}

			err = satellite.Overlay.Service.Put(ctx, n.Id, *n)
			assert.NoError(t, err, tt.name)
		}

		result, err := server.FindStorageNodes(ctx,
			&pb.FindStorageNodesRequest{
				Opts: &pb.OverlayOptions{Amount: tt.requestedNodeAmt},
			})

		if int64(tt.expectedResultLength) < tt.requestedNodeAmt {
			stat, ok := status.FromError(err)
			assert.Equal(t, true, ok, tt.name)
			assert.Equal(t, codes.ResourceExhausted, stat.Code(), tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}

		assert.Equal(t, tt.expectedResultLength, len(result.GetNodes()), tt.name)
	}
}
