package swarm_test

import (
	"context"
	"testing"

	swarmt "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-swarm/testing"

	peer "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-peer"
	transport "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-transport"
	ma "github.com/ipsn/go-ipfs/gxlibs/github.com/multiformats/go-multiaddr"
)

type dummyTransport struct {
	protocols []int
	proxy     bool
}

func (dt *dummyTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.Conn, error) {
	panic("unimplemented")
}

func (dt *dummyTransport) CanDial(addr ma.Multiaddr) bool {
	panic("unimplemented")
}

func (dt *dummyTransport) Listen(laddr ma.Multiaddr) (transport.Listener, error) {
	panic("unimplemented")
}

func (dt *dummyTransport) Proxy() bool {
	return dt.proxy
}

func (dt *dummyTransport) Protocols() []int {
	return dt.protocols
}

func TestUselessTransport(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	swarm := swarmt.GenSwarm(t, ctx)
	err := swarm.AddTransport(new(dummyTransport))
	if err == nil {
		t.Fatal("adding a transport that supports no protocols should have failed")
	}
}
