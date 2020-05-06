package model

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/satori/uuid"
)

// BlockData use to force a block data to meet certain func
type BlockData interface {
	toJSON() (string, error)
}

// Block is use to repesent block in blockchain
type Block struct {
	ID        uuid.UUID            `json:"id" gorm:"primary_key"`
	Timestamp time.Time            `json:"created_at"`
	Nonce     int32                `json:"nonce"`
	BlockType proto.BlockBlockType `json:"block_type"`
	PrevHash  string               `json:"prev_hash"`
	Data      string               `json:"data"`
	Hash      string               `json:"hash"`
}

// GenerateHash is used to generate based on block content
func (block *Block) GenerateHash() string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d:%d:%s:%s",
		block.ID.String(),
		block.Timestamp.Format(time.RFC3339),
		block.Nonce,
		block.BlockType,
		block.PrevHash,
		block.Data,
	)))
	return string(hash[:])
}

// GenerateNewBlock use to generate new block with nonce and
func GenerateNewBlock(
	id uuid.UUID,
	timestamp string,
	blockType proto.BlockBlockType,
	prevHash string,
	data BlockData,
) (*Block, error) {
	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return nil, err
	}

	jsondata, err := data.toJSON()
	if err != nil {
		return nil, err
	}

	var nonce int32 = 0
	newBlock := Block{
		id,
		t,
		nonce,
		blockType,
		prevHash,
		jsondata,
		"", // set to empty until a true hash generated
	}

	var hash string
	for hash = newBlock.GenerateHash(); len(hash) >= 3 && hash[:3] == "000"; hash = newBlock.GenerateHash() {
		nonce++
		newBlock = Block{
			id,
			t,
			nonce,
			blockType,
			prevHash,
			jsondata,
			"",
		}
	}

	return &Block{
		id,
		t,
		nonce,
		blockType,
		prevHash,
		jsondata,
		hash,
	}, nil
}
