package blockchain

import (
	"errors"
	"fmt"
	"sync"

	"github.com/neosouler7/bitGoin/db"
	"github.com/neosouler7/bitGoin/utils"
)

var (
	b           *blockchain
	once        sync.Once
	ErrNotFound = errors.New("block not found")
)

type blockchain struct {
	// now as a method of searching DB in a fast way
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1) // create & persist on db
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint() // search for checkpoint on db
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				fmt.Println("Blockchain Restoring...")
				b.restore(checkpoint) // restore b from bytes
			}
		})
	}
	return b
}
