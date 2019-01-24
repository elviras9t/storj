// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package bwagreement

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gtank/cryptopasta"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"

	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/peertls"
	"storj.io/storj/pkg/uplinkdb"
)

var (
	// Error the default bwagreement errs class
	Error = errs.Class("bwagreement error")
	mon   = monkit.Package()
)

// Config is a configuration struct that is everything you need to start an
// agreement receiver responsibility
type Config struct {
}

// DB stores bandwidth agreements.
type DB interface {
	// CreateAgreement adds a new bandwidth agreement.
	CreateAgreement(context.Context, string, Agreement) error
	// GetAgreements gets all bandwidth agreements.
	GetAgreements(context.Context) ([]Agreement, error)
	// GetAgreementsSince gets all bandwidth agreements since specific time.
	GetAgreementsSince(context.Context, time.Time) ([]Agreement, error)
}

// Server is an implementation of the pb.BandwidthServer interface
type Server struct {
	db       DB
	uplinkdb uplinkdb.DB
	pkey     crypto.PublicKey
	logger   *zap.Logger
}

// Agreement is a struct that contains a bandwidth agreement and the associated signature
type Agreement struct {
	Agreement []byte
	Signature []byte
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewServer creates instance of Server
func NewServer(db DB,  upldb uplinkdb.DB, logger *zap.Logger, pkey crypto.PublicKey) *Server {
	// TODO: reorder arguments, rename logger -> log
	return &Server{
		db:       db,
		uplinkdb: upldb,
		logger:   logger,
		pkey:     pkey,
	}
}

// Close closes resources
func (s *Server) Close() error { return nil }

// BandwidthAgreements receives and stores bandwidth agreements from storage nodes
func (s *Server) BandwidthAgreements(ctx context.Context, ba *pb.RenterBandwidthAllocation) (reply *pb.AgreementsSummary, err error) {
	defer mon.Task()(&ctx)(&err)

	s.logger.Debug("Received Agreement...")

	reply = &pb.AgreementsSummary{
		Status: pb.AgreementsSummary_REJECTED,
	}

	// storagenode signature is empty
	if len(ba.GetSignature()) == 0 {
		return reply, Error.New("Invalid Storage Node Signature length in the RenterBandwidthAllocation")
	}

	rbad := &pb.RenterBandwidthAllocation_Data{}
	if err = proto.Unmarshal(ba.GetData(), rbad); err != nil {
		return reply, Error.New("Failed to unmarshal RenterBandwidthAllocation: %+v", err)
	}

	pba := rbad.GetPayerAllocation()
	pbad := &pb.PayerBandwidthAllocation_Data{}
	if err := proto.Unmarshal(pba.GetData(), pbad); err != nil {
		return reply, Error.New("Failed to unmarshal PayerBandwidthAllocation: %+v", err)
	}

	// satellite signature is empty
	if len(pba.GetSignature()) == 0 {
		return reply, Error.New("Invalid Satellite Signature length in the PayerBandwidthAllocation")
	}

	if len(pbad.SerialNumber) == 0 {
		return reply, Error.New("Invalid SerialNumber in the PayerBandwidthAllocation")
	}

	if err = s.verifySignature(ctx, ba); err != nil {
		return reply, err
	}

	serialNum := pbad.GetSerialNumber() + rbad.StorageNodeId.String()

	// get and check expiration
	exp := time.Unix(pbad.GetExpirationUnixSec(), 0).UTC()
	if exp.Before(time.Now().UTC()) {
		return reply, Error.New("Bandwidth agreement is expired (%v)", exp)
	}

	err = s.db.CreateAgreement(ctx, serialNum, Agreement{
		Signature: ba.GetSignature(),
		Agreement: ba.GetData(),
		ExpiresAt: exp,
	})

	if err != nil {
		//todo:  better classify transport errors (AgreementsSummary_FAIL) vs logical (AgreementsSummary_REJECTED)
		return reply, Error.New("SerialNumber already exists in the PayerBandwidthAllocation")
	}

	reply.Status = pb.AgreementsSummary_OK
	s.logger.Debug("Stored Agreement...")
	return reply, nil
}

func (s *Server) verifySignature(ctx context.Context, ba *pb.RenterBandwidthAllocation) error {
	// TODO(security): detect replay attacks

	//Deserealize RenterBandwidthAllocation.GetData() so we can get public key
	rbad := &pb.RenterBandwidthAllocation_Data{}
	if err := proto.Unmarshal(ba.GetData(), rbad); err != nil {
		return Error.New("Failed to unmarshal RenterBandwidthAllocation: %+v", err)
	}

	pba := rbad.GetPayerAllocation()
	pbad := &pb.PayerBandwidthAllocation_Data{}
	if err := proto.Unmarshal(pba.GetData(), pbad); err != nil {
		return Error.New("Failed to unmarshal PayerBandwidthAllocation: %+v", err)
	}

	// Get renter's public key from uplink agreement db
	uplinkInfo, err := s.uplinkdb.GetSignature(ctx, pbad.GetSerialNumber())
	if err != nil {
		return BwAgreementError.New("Failed to unmarshal PayerBandwidthAllocation: %+v", err)
	}

	pubkey, err := x509.ParsePKIXPublicKey(uplinkInfo.Signature)
	if err != nil {
		return Error.New("Failed to extract Public Key from RenterBandwidthAllocation: %+v", err)
	}

	// Typecast public key
	k, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		return peertls.ErrUnsupportedKey.New("%T", pubkey)
	}

	signatureLength := k.Curve.Params().P.BitLen() / 8
	if len(ba.GetSignature()) < signatureLength {
		return Error.New("Invalid Renter's Signature Length")
	}
	// verify Renter's (uplink) signature
	if ok := cryptopasta.Verify(ba.GetData(), ba.GetSignature(), k); !ok {
		return Error.New("Failed to verify Renter's Signature")
	}

	// satellite public key
	k, ok = s.pkey.(*ecdsa.PublicKey)
	if !ok {
		return peertls.ErrUnsupportedKey.New("%T", s.pkey)
	}

	signatureLength = k.Curve.Params().P.BitLen() / 8
	if len(rbad.GetPayerAllocation().GetSignature()) < signatureLength {
		return Error.New("Inavalid Payer's Signature Length")
	}
	// verify Payer's (satellite) signature
	if ok := cryptopasta.Verify(rbad.GetPayerAllocation().GetData(), rbad.GetPayerAllocation().GetSignature(), k); !ok {
		return Error.New("Failed to verify Payer's Signature")
	}
	return nil
}
