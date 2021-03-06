package commands

import (
	"testing"

	"github.com/ipsn/go-ipfs/namesys"

	tu "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-testutil"
	ipns "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipns"
)

func TestKeyTranslation(t *testing.T) {
	pid := tu.RandPeerIDFatal(t)
	pkname := namesys.PkKeyForID(pid)
	ipnsname := ipns.RecordKey(pid)

	pkk, err := escapeDhtKey("/pk/" + pid.Pretty())
	if err != nil {
		t.Fatal(err)
	}

	ipnsk, err := escapeDhtKey("/ipns/" + pid.Pretty())
	if err != nil {
		t.Fatal(err)
	}

	if pkk != pkname {
		t.Fatal("keys didnt match!")
	}

	if ipnsk != ipnsname {
		t.Fatal("keys didnt match!")
	}
}
