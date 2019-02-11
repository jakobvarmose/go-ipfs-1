package cbstore

import (
	"fmt"

	blocks "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-block-format"
	cid "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-cid"
	bstore "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-merkledag"
)

type blockstore struct {
	bstore.Blockstore
}

type publicInterface interface {
	Public() (blocks.Block, error)
}

func NewBlockstore(bs bstore.Blockstore) bstore.Blockstore {
	return &blockstore{bs}
}

func (bs *blockstore) DeleteBlock(c cid.Cid) error {
	if c.EncryptionAlgorithm() == 0 {
		return bs.Blockstore.DeleteBlock(c)
	}

	return bs.Blockstore.DeleteBlock(c.Public())
}

func (bs *blockstore) Has(c cid.Cid) (bool, error) {
	if c.EncryptionAlgorithm() == 0 {
		return bs.Blockstore.Has(c)
	}

	return bs.Blockstore.Has(c.Public())
}

func (bs *blockstore) Get(c cid.Cid) (blocks.Block, error) {
	fmt.Println("GET", c, c.Public())
	if c.EncryptionAlgorithm() == 0 {
		return bs.Blockstore.Get(c)
	}

	block, err := bs.Blockstore.Get(c.Public())
	if err != nil {
		return nil, err
	}

	return blocks.NewAutoDecryptedBlock(block.RawData(), c)
}

func (bs *blockstore) GetSize(c cid.Cid) (int, error) {
	if c.EncryptionAlgorithm() == 0 {
		return bs.Blockstore.GetSize(c)
	}

	block, err := bs.Blockstore.Get(c.Public())
	if err != nil {
		return 0, err
	}

	block, err = blocks.NewAutoDecryptedBlock(block.RawData(), c)
	if err != nil {
		return 0, err
	}

	return len(block.RawData()), nil
}

func (bs *blockstore) Put(block blocks.Block) error {
	fmt.Println("PUT", block.Cid(), block.Cid().Public())
	fmt.Printf("%T %v\n", block, block)
	x, ok := block.(*merkledag.ProtoNode)
	if ok {
		fmt.Println("OK", len(x.Links()), x.Links())
	}
	block2, ok := block.(publicInterface)
	if !ok {
		return bs.Blockstore.Put(block)
	}

	block, err := block2.Public()
	if err != nil {
		return err
	}
	fmt.Println("XXX", block.Cid())

	return bs.Blockstore.Put(block)
}

func (bs *blockstore) PutMany(b1 []blocks.Block) error {
	b2 := make([]blocks.Block, len(b1))
	for i, block := range b1 {
		block2, ok := block.(publicInterface)
		if !ok {
			b2[i] = block
			continue
		}

		block, err := block2.Public()
		if err != nil {
			return err
		}
		b2[i] = block
	}

	return bs.Blockstore.PutMany(b2)
}
