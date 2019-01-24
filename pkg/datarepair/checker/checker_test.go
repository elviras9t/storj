// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

//TODO: reenable
// +build ignore

package checker_test

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/storj/internal/testcontext"
	"storj.io/storj/internal/testplanet"
	"storj.io/storj/internal/teststorj"
	"storj.io/storj/pkg/auth"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/satellite/satellitedb"
)

func TestIdentifyInjuredSegments(t *testing.T) {
	t.Skip("needs update")

	// logic should be roughly:
	// use storagenodes as the valid and
	// generate invalid ids
	// identify should then find the invalid ones
	// note satellite's: own sub-systems need to be disabled

	tctx := testcontext.New(t)
	defer tctx.Cleanup()

	planet, err := testplanet.New(t, 1, 4, 0)
	require.NoError(t, err)
	defer tctx.Check(planet.Shutdown)

	planet.Start(tctx)
	time.Sleep(2 * time.Second)

	pointerdb := planet.Satellites[0].Metainfo.Endpoint
	repairQueue := planet.Satellites[0].DB.RepairQueue()

	ctx := auth.WithAPIKey(tctx, nil)

	const N = 25
	nodes := []*pb.Node{}
	segs := []*pb.InjuredSegment{}
	//fill a pointerdb
	for i := 0; i < N; i++ {
		s := strconv.Itoa(i)
		ids := teststorj.NodeIDsFromStrings([]string{s + "a", s + "b", s + "c", s + "d"}...)

		p := &pb.Pointer{
			Remote: &pb.RemoteSegment{
				Redundancy: &pb.RedundancyScheme{
					RepairThreshold: int32(2),
				},
				PieceId: strconv.Itoa(i),
				RemotePieces: []*pb.RemotePiece{
					{PieceNum: 0, NodeId: ids[0]},
					{PieceNum: 1, NodeId: ids[1]},
					{PieceNum: 2, NodeId: ids[2]},
					{PieceNum: 3, NodeId: ids[3]},
				},
			},
		}

		req := &pb.PutRequest{
			Path:    p.Remote.PieceId,
			Pointer: p,
		}

		resp, err := pointerdb.Put(ctx, req)
		assert.NotNil(t, resp)
		assert.NoError(t, err)

		//nodes for cache
		selection := rand.Intn(4)
		for _, v := range ids[:selection] {
			n := &pb.Node{Id: v, Type: pb.NodeType_STORAGE, Address: &pb.NodeAddress{Address: ""}}
			nodes = append(nodes, n)
		}

		pieces := []int32{0, 1, 2, 3}
		//expected injured segments
		if len(ids[:selection]) < int(p.Remote.Redundancy.RepairThreshold) {
			seg := &pb.InjuredSegment{
				Path:       p.Remote.PieceId,
				LostPieces: pieces[selection:],
			}
			segs = append(segs, seg)
		}
	}

	checker := planet.Satellites[0].Repair.Checker
	assert.NoError(t, err)
	err = checker.IdentifyInjuredSegments(ctx)
	assert.NoError(t, err)

	expected := map[string]*pb.InjuredSegment{}
	for _, seg := range segs {
		expected[seg.Path] = seg
	}

	//check if the expected segments were added to the queue
	dequeued := []*pb.InjuredSegment{}
	for i := 0; i < len(segs); i++ {
		injSeg, err := repairQueue.Dequeue(ctx)
		assert.NoError(t, err)

		if _, ok := expected[injSeg.Path]; ok {
			t.Log("got", injSeg.Path)
			delete(expected, injSeg.Path)
		} else {
			t.Error("unexpected", injSeg)
		}
		dequeued = append(dequeued, &injSeg)
	}

	for _, missing := range expected {
		t.Error("did not get", missing)
	}
}

func TestOfflineNodes(t *testing.T) {
	t.Skip("needs update")

	tctx := testcontext.New(t)
	defer tctx.Cleanup()

	planet, err := testplanet.New(t, 1, 0, 0)
	require.NoError(t, err)
	defer tctx.Check(planet.Shutdown)

	planet.Start(tctx)
	time.Sleep(2 * time.Second)

	const N = 50
	nodes := []*pb.Node{}
	nodeIDs := storj.NodeIDList{}
	expectedOffline := []int32{}
	for i := 0; i < N; i++ {
		id := teststorj.NodeIDFromString(strconv.Itoa(i))
		n := &pb.Node{Id: id, Type: pb.NodeType_STORAGE, Address: &pb.NodeAddress{Address: ""}}
		nodes = append(nodes, n)
		if i%(rand.Intn(5)+2) == 0 {
			nodeIDs = append(nodeIDs, teststorj.NodeIDFromString("id"+id.String()))
			expectedOffline = append(expectedOffline, int32(i))
		} else {
			nodeIDs = append(nodeIDs, id)
		}
	}

	checker := planet.Satellites[0].Repair.Checker
	assert.NoError(t, err)
	offline, err := checker.OfflineNodes(tctx, nodeIDs)
	assert.NoError(t, err)
	assert.Equal(t, expectedOffline, offline)
}

func BenchmarkIdentifyInjuredSegments(b *testing.B) {
	b.Skip("needs update")

	tctx := testcontext.New(b)
	defer tctx.Cleanup()

	planet, err := testplanet.New(b, 1, 0, 0)
	require.NoError(b, err)
	defer tctx.Check(planet.Shutdown)

	planet.Start(tctx)
	time.Sleep(2 * time.Second)

	pointerdb := planet.Satellites[0].Metainfo.Endpoint
	repairQueue := planet.Satellites[0].DB.RepairQueue()

	ctx := auth.WithAPIKey(tctx, nil)

	// creating in-memory db and opening connection
	db, err := satellitedb.NewInMemory()
	defer func() {
		err = db.Close()
		assert.NoError(b, err)
	}()
	err = db.CreateTables()
	assert.NoError(b, err)

	const N = 25
	nodes := []*pb.Node{}
	segs := []*pb.InjuredSegment{}
	//fill a pointerdb
	for i := 0; i < N; i++ {
		s := strconv.Itoa(i)
		ids := teststorj.NodeIDsFromStrings([]string{s + "a", s + "b", s + "c", s + "d"}...)

		p := &pb.Pointer{
			Remote: &pb.RemoteSegment{
				Redundancy: &pb.RedundancyScheme{
					RepairThreshold: int32(2),
				},
				PieceId: strconv.Itoa(i),
				RemotePieces: []*pb.RemotePiece{
					{PieceNum: 0, NodeId: ids[0]},
					{PieceNum: 1, NodeId: ids[1]},
					{PieceNum: 2, NodeId: ids[2]},
					{PieceNum: 3, NodeId: ids[3]},
				},
			},
		}
		req := &pb.PutRequest{
			Path:    p.Remote.PieceId,
			Pointer: p,
		}

		resp, err := pointerdb.Put(ctx, req)
		assert.NotNil(b, resp)
		assert.NoError(b, err)

		//nodes for cache
		selection := rand.Intn(4)
		for _, v := range ids[:selection] {
			n := &pb.Node{Id: v, Type: pb.NodeType_STORAGE, Address: &pb.NodeAddress{Address: ""}}
			nodes = append(nodes, n)
		}
		pieces := []int32{0, 1, 2, 3}
		//expected injured segments
		if len(ids[:selection]) < int(p.Remote.Redundancy.RepairThreshold) {
			seg := &pb.InjuredSegment{
				Path:       p.Remote.PieceId,
				LostPieces: pieces[selection:],
			}
			segs = append(segs, seg)
		}
	}
	//fill a overlay cache
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker := planet.Satellites[0].Repair.Checker
		assert.NoError(b, err)

		err = checker.IdentifyInjuredSegments(ctx)
		assert.NoError(b, err)

		//check if the expected segments were added to the queue
		dequeued := []*pb.InjuredSegment{}
		for i := 0; i < len(segs); i++ {
			injSeg, err := repairQueue.Dequeue(ctx)
			assert.NoError(b, err)
			dequeued = append(dequeued, &injSeg)
		}
		sort.Slice(segs, func(i, k int) bool { return segs[i].Path < segs[k].Path })
		sort.Slice(dequeued, func(i, k int) bool { return dequeued[i].Path < dequeued[k].Path })

		for i := 0; i < len(segs); i++ {
			assert.True(b, proto.Equal(segs[i], dequeued[i]))
		}
	}
}
