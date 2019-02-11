package cid

import (
	mh "github.com/ipsn/go-ipfs/gxlibs/github.com/multiformats/go-multihash"
	"github.com/ipsn/go-ipfs/multisymmetric"
)

type Builder interface {
	Sum(data []byte) (Cid, error)
	GetCodec() uint64
	WithCodec(uint64) Builder
}

type V0Builder struct{}

type V1Builder struct {
	Codec    uint64
	MhType   uint64
	MhLength int // MhLength <= 0 means the default length
}

func (p Prefix) GetCodec() uint64 {
	return p.Codec
}

func (p Prefix) WithCodec(c uint64) Builder {
	if c == p.Codec {
		return p
	}
	p.Codec = c
	if c != DagProtobuf {
		p.Version = 1
	}
	return p
}

func (p V0Builder) Sum(data []byte) (Cid, error) {
	hash, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		return Undef, err
	}
	return NewCidV0(hash), nil
}

func (p V0Builder) GetCodec() uint64 {
	return DagProtobuf
}

func (p V0Builder) WithCodec(c uint64) Builder {
	if c == DagProtobuf {
		return p
	}
	return V1Builder{Codec: c, MhType: mh.SHA2_256}
}

func (p V1Builder) Sum(data []byte) (Cid, error) {
	mhLen := p.MhLength
	if mhLen <= 0 {
		mhLen = -1
	}
	hash, err := mh.Sum(data, p.MhType, mhLen)
	if err != nil {
		return Undef, err
	}
	return NewCidV1(p.Codec, hash), nil
}

func (p V1Builder) GetCodec() uint64 {
	return p.Codec
}

func (p V1Builder) WithCodec(c uint64) Builder {
	p.Codec = c
	return p
}

type V4Builder struct {
	Codec               uint64
	MhType              uint64
	MhLength            int // MhLength <= 0 means the default length
	EncryptionAlgorithm uint64
}

func (p V4Builder) Sum(data []byte) (Cid, error) {
	key, err := multisymmetric.GenerateKey(p.EncryptionAlgorithm)
	if err != nil {
		return Undef, err
	}
	ct, err := multisymmetric.Encrypt(p.EncryptionAlgorithm, key, data)
	if err != nil {
		return Undef, err
	}
	hash, err := mh.Sum(ct, p.MhType, p.MhLength)
	if err != nil {
		return Undef, err
	}
	c := NewCidV1(Raw, hash)
	return NewPrivateCid(p.Codec, p.EncryptionAlgorithm, key, c), nil
}

func (p V4Builder) GetCodec() uint64 {
	return p.Codec
}

func (p V4Builder) WithCodec(c uint64) Builder {
	p.Codec = c
	return p
}
