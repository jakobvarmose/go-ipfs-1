package cidutil

import (
	cid "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-cid"
	mhash "github.com/ipsn/go-ipfs/gxlibs/github.com/multiformats/go-multihash"
	"github.com/ipsn/go-ipfs/multisymmetric"
)

// InlineBuilder is a cid.Builder that will use the id multihash when the
// size of the content is no more than limit
type InlineBuilder struct {
	cid.Builder     // Parent Builder
	Limit       int // Limit (inclusive)
}

// WithCodec implements the cid.Builder interface
func (p InlineBuilder) WithCodec(c uint64) cid.Builder {
	return InlineBuilder{p.Builder.WithCodec(c), p.Limit}
}

// Sum implements the cid.Builder interface
func (p InlineBuilder) Sum(data []byte) (cid.Cid, error) {
	if len(data) > p.Limit {
		return p.Builder.Sum(data)
	}
	return cid.V4Builder{Codec: p.GetCodec(), MhType: mhash.ID, EncryptionAlgorithm: multisymmetric.AES_CTR_ZERO_IV}.Sum(data)
}
