package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookCheckout struct {
	BookID    string `json:"book_id"`
	User      string `json:"user"`
	IsGenesis bool   `json:"-"`
}

type Blockchain struct {
	Blocks []Block
}

type Block struct {
	PrevHash  string       `json:"prev_hash"`
	Hash      string       `json:"hash"`
	Data      BookCheckout `json:"data"`
	Timestamp string       `json:"timestamp"`
	Position  int          `json:"position"`
}

func createGenesisBlock() *Block {
  return createBlock(&Block{}, &BookCheckout{IsGenesis: true})
}

func NewBlockChain() *Blockchain {
	return &Blockchain{
    Blocks: []Block{*createGenesisBlock()},
  }
}

func createBlock(prevBlock *Block, data *BookCheckout) *Block {
	block := &Block{
		PrevHash:  prevBlock.Hash,
		Data:      *data,
		Position:  prevBlock.Position + 1,
		Timestamp: time.Now().String(),
	}
	if data.IsGenesis {
		block.PrevHash = ""
    block.Position = 0
	}
  block.generateHash()
	return block
}

func (bc *Blockchain) WriteBlock(data *BookCheckout) error {
  prevBlock := &bc.Blocks[len(bc.Blocks)-1]
  newBlock := createBlock(prevBlock, data)
  
  if err := bc.ValidateBlock(newBlock); err != nil {
    return err 
  }

  bc.Blocks = append(bc.Blocks, *newBlock)
  return nil
}

func (bc *Blockchain) ValidateBlock(block *Block) error {
  prevBlock := &bc.Blocks[len(bc.Blocks)-1]
  if prevBlock.Hash != block.PrevHash {
    return fmt.Errorf("Previous hash not equal")
  } 
  if prevBlock.Position + 1 != block.Position {
    return fmt.Errorf("Position not correct")
  }
  if !block.validateHash(block.Hash) {
    return fmt.Errorf("Hash not correct")
  }
  return nil
}

func (block *Block) generateHash() {
  hash := sha256.New()
  stringData, _ := json.MarshalIndent(block.Data, "", "  ")
  data := block.PrevHash + block.Timestamp + string(block.Position) + string(stringData)

  io.WriteString(hash, data)
  block.Hash = string(hash.Sum(nil))
}

func (block *Block) validateHash(hash string) bool {
  block.generateHash()
  if block.Hash != hash {
    return false
  }
  return true
}
