// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package pointerdb

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/skyrings/skyring-common/tools/uuid"

	"storj.io/storj/pkg/auth"
	"storj.io/storj/pkg/identity"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/peertls"
	"storj.io/storj/pkg/uplinkdb"
)

// AllocationSigner structure
type AllocationSigner struct {
	satelliteIdentity *identity.FullIdentity
	bwExpiration      int
	uplinkdb          uplinkdb.DB
}

// NewAllocationSigner creates new instance
func NewAllocationSigner(satelliteIdentity *identity.FullIdentity, bwExpiration int, upldb uplinkdb.DB) *AllocationSigner {
	return &AllocationSigner{
		satelliteIdentity: satelliteIdentity,
		bwExpiration:      bwExpiration,
		uplinkdb:          upldb,
	}
}

// PayerBandwidthAllocation returns generated payer bandwidth allocation
func (allocation *AllocationSigner) PayerBandwidthAllocation(ctx context.Context, peerIdentity *identity.PeerIdentity, action pb.PayerBandwidthAllocation_Action) (pba *pb.PayerBandwidthAllocation, err error) {
	if peerIdentity == nil {
		return nil, Error.New("missing peer identity")
	}

	pk, ok := peerIdentity.Leaf.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, peertls.ErrUnsupportedKey.New("%T", peerIdentity.Leaf.PublicKey)
	}

	pubbytes, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return nil, err
	}

	created := time.Now().Unix()

	// convert ttl from days to seconds
	ttl := allocation.bwExpiration
	ttl *= 86400

	serialNum, err := uuid.New()
	if err != nil {
		return nil, err
	}

	// store the corresponding uplink's id and public key into uplinkDB db
	err = allocation.uplinkdb.CreateAgreement(ctx, serialNum.String(), uplinkdb.Agreement{Agreement: peerIdentity.ID.Bytes(), Signature: pubbytes})
	if err != nil {
		return nil, err
	}

	pbad := &pb.PayerBandwidthAllocation_Data{
		SatelliteId:       allocation.satelliteIdentity.ID,
		UplinkId:          peerIdentity.ID,
		CreatedUnixSec:    created,
		ExpirationUnixSec: created + int64(ttl),
		Action:            action,
		SerialNumber:      serialNum.String(),
		PubKey:            pubbytes,
	}

	data, err := proto.Marshal(pbad)
	if err != nil {
		return nil, err
	}
	signature, err := auth.GenerateSignature(data, allocation.satelliteIdentity)
	if err != nil {
		return nil, err
	}
	return &pb.PayerBandwidthAllocation{Signature: signature, Data: data}, nil
}
