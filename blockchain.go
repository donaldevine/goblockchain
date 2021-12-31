package main

import (
        "crypto/sha256"
        "encoding/json"
        "fmt"
        "strconv"
        "strings"
        "time"
)

type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

// calculate hash function
func (b Block) calculateHash() string {
	//convert block data to json
	data, _ := json.Marshal(b.data)
	
	//concatenate block's previous hash, data, timestamp and proof of work
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)

	//hash it all with sha 256	
	blockHash := sha256.Sum256([]byte(blockData))

	//return it as base 16
	return fmt.Sprintf("%x", blockHash)
}


//mining function
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
			b.pow++
			b.hash = b.calculateHash()
	}
}

//create initial block
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
			hash:      "0",
			timestamp: time.Now(),
	}
	return Blockchain{
			genesisBlock,
			[]Block{genesisBlock},
			difficulty,
	}
}

//this function adds a new block
func (b *Blockchain) addBlock(from, to string, amount float64) {

	//create a block data with the data
	blockData := map[string]interface{}{
			"from":   from,
			"to":     to,
			"amount": amount,
	}

	//get the last block on the chain
	lastBlock := b.chain[len(b.chain)-1]

	//create a new block with the data and include the previous block hash
	newBlock := Block{
			data:         blockData,
			previousHash: lastBlock.hash,
			timestamp:    time.Now(),
	}

	//mine this new block
	newBlock.mine(b.difficulty)

	//append it to the chain
	b.chain = append(b.chain, newBlock)
}

//check if a blockchain is valid
func (b Blockchain) isValid() bool {
	//for each block in the chain we check that the hash matches the calculated hash 
	//and that the previous block hash matches the has of the previous hash property of the current block
	for i := range b.chain[1:] {
			previousBlock := b.chain[i]
			currentBlock := b.chain[i+1]
			if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
					return false
			}
	}
	return true
}

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(2)

	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.addBlock("Alice", "Bob", 5)
	blockchain.addBlock("John", "Bob", 2)

	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid())
}