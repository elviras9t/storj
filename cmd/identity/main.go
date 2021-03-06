// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"storj.io/storj/internal/fpath"
	"storj.io/storj/pkg/certificates"
	"storj.io/storj/pkg/cfgstruct"
	"storj.io/storj/pkg/identity"
	"storj.io/storj/pkg/peertls"
	"storj.io/storj/pkg/process"
)

const (
	defaultSignerAddress = "certs.alpha.storj.io:8888"
)

var (
	rootCmd = &cobra.Command{
		Use:   "identity",
		Short: "Identity management",
	}

	newServiceCmd = &cobra.Command{
		Use:         "create <service>",
		Short:       "Create a new full identity for a service",
		Args:        cobra.ExactArgs(1),
		RunE:        cmdNewService,
		Annotations: map[string]string{"type": "setup"},
	}

	csrCmd = &cobra.Command{
		Use:         "authorize <service> <auth-token>",
		Short:       "Send a certificate signing request for a service's CA certificate",
		Args:        cobra.ExactArgs(2),
		RunE:        cmdCSR,
		Annotations: map[string]string{"type": "setup"},
	}

	//nolint
	config struct {
		Difficulty     uint64 `default:"30" help:"minimum difficulty for identity generation"`
		Concurrency    uint   `default:"4" help:"number of concurrent workers for certificate authority generation"`
		ParentCertPath string `help:"path to the parent authority's certificate chain"`
		ParentKeyPath  string `help:"path to the parent authority's private key"`
		Signer         certificates.CertClientConfig
	}

	identityDir        string
	defaultIdentityDir = fpath.ApplicationDir("storj", "identity")
)

func init() {
	rootCmd.PersistentFlags().StringVar(&identityDir, "identity-dir", defaultIdentityDir, "root directory for identity output")

	rootCmd.AddCommand(newServiceCmd)
	rootCmd.AddCommand(csrCmd)

	cfgstruct.Bind(newServiceCmd.Flags(), &config, cfgstruct.IdentityDir(defaultIdentityDir))
	cfgstruct.Bind(csrCmd.Flags(), &config, cfgstruct.IdentityDir(defaultIdentityDir))
}

func main() {
	process.Exec(rootCmd)
}

func serviceDirectory(serviceName string) string {
	return filepath.Join(identityDir, serviceName)
}

func cmdNewService(cmd *cobra.Command, args []string) error {
	serviceDir := serviceDirectory(args[0])

	caCertPath := filepath.Join(serviceDir, "ca.cert")
	caKeyPath := filepath.Join(serviceDir, "ca.key")
	identCertPath := filepath.Join(serviceDir, "identity.cert")
	identKeyPath := filepath.Join(serviceDir, "identity.key")

	caConfig := identity.CASetupConfig{
		CertPath:       caCertPath,
		KeyPath:        caKeyPath,
		Difficulty:     config.Difficulty,
		Concurrency:    config.Concurrency,
		ParentCertPath: config.ParentCertPath,
		ParentKeyPath:  config.ParentKeyPath,
	}

	if caConfig.Status() != identity.NoCertNoKey {
		return errs.New("CA certificate and/or key already exits, NOT overwriting!")
	}

	ca, caerr := caConfig.Create(process.Ctx(cmd))
	if caerr != nil {
		return caerr
	}

	identConfig := identity.SetupConfig{
		CertPath: identCertPath,
		KeyPath:  identKeyPath,
	}

	if identConfig.Status() != identity.NoCertNoKey {
		return errs.New("Identity certificate and/or key already exits, NOT overwriting!")
	}

	_, iderr := identConfig.Create(ca)
	if iderr != nil {
		return iderr
	}

	fmt.Printf("Unsigned identity is located in %q\n", serviceDir)
	return nil
}

func cmdCSR(cmd *cobra.Command, args []string) error {
	ctx := process.Ctx(cmd)

	serviceDir := serviceDirectory(args[0])
	authToken := args[1]

	caCertPath := filepath.Join(serviceDir, "ca.cert")
	caKeyPath := filepath.Join(serviceDir, "ca.key")
	caConfig := identity.FullCAConfig{
		CertPath: caCertPath,
		KeyPath:  caKeyPath,
	}
	identCertPath := filepath.Join(serviceDir, "identity.cert")
	identKeyPath := filepath.Join(serviceDir, "identity.key")
	identConfig := identity.Config{
		CertPath: identCertPath,
		KeyPath:  identKeyPath,
	}

	ca, err := caConfig.Load()
	if err != nil {
		return err
	}
	ident, err := identConfig.Load()
	if err != nil {
		return err
	}

	if config.Signer.Address == "" {
		config.Signer.Address = defaultSignerAddress
	}

	signedChainBytes, err := config.Signer.Sign(ctx, ident, authToken)
	if err != nil {
		return errs.New("error occurred while signing certificate: %s\n(identity files were still generated and saved, if you try again existing files will be loaded)", err)
	}

	signedChain, err := identity.ParseCertChain(signedChainBytes)
	if err != nil {
		return nil
	}

	err = caConfig.SaveBackup(ca)
	if err != nil {
		return err
	}

	ca.Cert = signedChain[0]
	ca.RestChain = signedChain[1:]
	err = identity.FullCAConfig{
		CertPath: caConfig.CertPath,
	}.Save(ca)
	if err != nil {
		return err
	}

	err = identConfig.SaveBackup(ident)
	if err != nil {
		return err
	}

	ident.RestChain = signedChain[1:]
	ident.CA = ca.Cert
	err = identity.Config{
		CertPath: identConfig.CertPath,
	}.Save(ident)
	if err != nil {
		return err
	}

	fmt.Printf("Signed identity is in %q\n", serviceDir)
	return nil
}

func printExtensions(cert []byte, exts []pkix.Extension) error {
	hash, err := peertls.SHA256Hash(cert)
	if err != nil {
		return err
	}
	b64Hash, err := json.Marshal(hash)
	if err != nil {
		return err
	}
	fmt.Printf("Cert hash: %s\n", b64Hash)
	fmt.Println("Extensions:")
	for _, e := range exts {
		var data interface{}
		switch e.Id.String() {
		case peertls.ExtensionIDs[peertls.RevocationExtID].String():
			var rev peertls.Revocation
			if err := rev.Unmarshal(e.Value); err != nil {
				return err
			}
			data = rev
		default:
			data = e.Value
		}
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("\t%s: %s\n", e.Id, out)
	}
	return nil
}
