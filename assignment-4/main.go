package main

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
)

const (
	TargetURL = "http://crypto-class.appspot.com/po?er="
)

func Do(s string) bool {
	escaped := url.QueryEscape(s)

	resp, err := http.Post(TargetURL+escaped, "application/text", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Printf("got response code: %v\n", resp.StatusCode)
	if resp.StatusCode == http.StatusNotFound {
		return true // good padding
	}
	return false // bad padding
}

type BlockPair struct {
	first  []byte
	second []byte
	offset int
}

func main() {
	// aes cbc encrypted cipherthext with IV as the first block.
	ciphertext := mustBytes(hex.DecodeString("f20bdba6ff29eed7b046d1df9fb7000058b1ffb4210a580f748b4ac714c001bd4a61044426fb515dad3f21f18aa577c0bdf302936266926ff37dbf7035d5eeb4"))

	for _, pair := range blockPairs(ciphertext, aes.BlockSize)[2:] {
		var guessed []byte

		for padding := 1; padding <= aes.BlockSize; padding++ {
			modified := bytes.Clone(pair.first)
			for i := 1; i <= padding-1; i++ {
				modified[len(modified)-i] ^= guessed[i-1] ^ byte(padding)
			}

			for guess := 2; guess < 256; guess++ {
				offset := len(modified) - 1 - len(guessed)
				modified[offset] = modified[offset] ^ byte(guess) ^ byte(padding)
				if Do(hex.EncodeToString(append(modified, pair.second...))) {
					guessed = append(guessed, byte(guess))
					break
				}
				modified[offset] ^= byte(guess) ^ byte(padding)
			}
		}

		fmt.Println(string(rev(guessed)))
	}
}

func blockPairs(msg []byte, blockSize int) []BlockPair {
	var result []BlockPair

	for block := 1; block < len(msg)/blockSize; block++ {
		result = append(result, BlockPair{
			first:  bytes.Clone(msg[(block-1)*blockSize : block*blockSize]),
			second: bytes.Clone(msg[block*blockSize : (block+1)*blockSize]),
			offset: block,
		})
	}

	return result
}

func mustBytes(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}

func rev(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-1-i] = b[len(b)-1-i], b[i]
	}
	return b
}
