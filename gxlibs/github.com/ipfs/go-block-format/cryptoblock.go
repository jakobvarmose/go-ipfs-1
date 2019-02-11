package blocks

import (
	"fmt"

	cid "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-cid"
	mh "github.com/ipsn/go-ipfs/gxlibs/github.com/multiformats/go-multihash"
	"github.com/ipsn/go-ipfs/multisymmetric"
)

type CryptoBlock struct {
	cid        cid.Cid
	plaintext  []byte
	ciphertext []byte
}

// NewBlockWithPrefix creates a new block from data and a prefix.
// If the prefix is a version 4 prefix the data will be encrypted.
func NewBlockWithPrefix(plaintext []byte, pref cid.Prefix) (Block, error) {
	if pref.Version != 4 {
		c, err := pref.Sum(plaintext)
		if err != nil {
			return nil, err
		}
		return NewBlockWithCid(plaintext, c)
	}

	key, err := multisymmetric.GenerateKey(pref.EncryptionAlgorithm)
	if err != nil {
		return nil, err
	}

	ciphertext, err := multisymmetric.Encrypt(pref.EncryptionAlgorithm, key, plaintext)
	if err != nil {
		return nil, err
	}

	cid1, err := cid.NewPrefixV1(cid.Raw, pref.MhType).Sum(ciphertext)
	if err != nil {
		return nil, err
	}

	cid4 := cid.NewPrivateCid(pref.Codec, pref.EncryptionAlgorithm, key, cid1)

	return &CryptoBlock{
		cid:        cid4,
		plaintext:  plaintext,
		ciphertext: ciphertext,
	}, nil
}

// NewAutoDecryptedBlock creates a new block from encrypted data and a cid.
func NewAutoDecryptedBlock(ciphertext []byte, c cid.Cid) (Block, error) {
	if c.Prefix().Version != 4 {
		return NewBlockWithCid(ciphertext, c)
	}

	plaintext, err := multisymmetric.Decrypt(c.EncryptionAlgorithm(), c.EncryptionKey(), ciphertext)
	if err != nil {
		return nil, err
	}

	return &CryptoBlock{
		cid:        c,
		plaintext:  plaintext,
		ciphertext: ciphertext,
	}, nil
}

// Multihash returns the hash contained in the block CID.
func (b *CryptoBlock) Multihash() mh.Multihash {
	return b.cid.Hash()
}

// RawData returns the block raw contents as a byte slice.
func (b *CryptoBlock) RawData() []byte {
	return b.plaintext
}

// Cid returns the content identifier of the block.
func (b *CryptoBlock) Cid() cid.Cid {
	return b.cid
}

// String provides a human-readable representation of the block CID.
func (b *CryptoBlock) String() string {
	return fmt.Sprintf("[CryptoBlock %s]", b.Cid())
}

// Loggable returns a go-log loggable item.
func (b *CryptoBlock) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"cryptoBlock": b.Cid().String(),
	}
}

func (b *CryptoBlock) Public() (Block, error) {
	return NewBlockWithCid(b.ciphertext, b.cid.Public())
}
