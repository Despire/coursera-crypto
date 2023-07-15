package main

import (
	"crypto/sha256"
	"fmt"
	"math"
	"os"
)

const (
    filePath = "./6.1.intro.mp4_download"
    blockSize = 1 << 10 // 1KB
)

func main() {
    contents, err := os.ReadFile(filePath)
    if err != nil {
        panic(err)
    }

    blocks := splitIntoBlocks(contents)
    hash := sha256.Sum256(blocks[len(blocks) - 1].block)

    for i := len(blocks) - 2; i >= 0; i-- {
        hash = sha256.Sum256(append(blocks[i].block, hash[:]...))
    }

    fmt.Printf("%x\n", hash)
}

type Block struct { block []byte }

func splitIntoBlocks(b []byte) []Block {
    blockLength := len(b) / blockSize
    if len(b) % blockSize != 0 {
        blockLength++
    }

    result := make([]Block, 0, blockLength)

    for i := 0; i < blockLength; i++ {
        result = append(result, Block{
            block: b[i*blockSize: int(math.Min(float64(len(b)), float64(i*blockSize+blockSize)))],
        })
    }

    return result
}
