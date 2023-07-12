package main

import (
	"crypto/aes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

var (
	// ErrInvalidPKCS5Padding is returned when a message has invalid pkcs5 padding.
	ErrInvalidPKCS5Padding = errors.New("invalid pkcs5 padded message")
	// ErrInvalidBlockSize is returned when the message is not of the multiple of the block size of the cipher.
	ErrInvalidBlockSize = errors.New("invalid block size")
)

func pkcs5Validate(block []byte) ([]byte, error) {
	paddedBytes := block[len(block)-1]
	for i := len(block) - 1; i > len(block)-int(paddedBytes); i-- {
		if block[i] != paddedBytes {
			return nil, ErrInvalidPKCS5Padding
		}
	}

	return block[:len(block)-int(paddedBytes)], nil
}

func aesDecryptCBC(key, msg []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(msg)%aes.BlockSize != 0 {
		return nil, ErrInvalidBlockSize
	}

	prevBlock := msg[:aes.BlockSize]
	for i := aes.BlockSize; i < len(msg); i += aes.BlockSize {
		ciphertext := append([]byte(nil), msg[i:i+aes.BlockSize]...)
		current := msg[i : i+aes.BlockSize]
		cipher.Decrypt(current, current)
		copy(current, xor(current, prevBlock))
		prevBlock = ciphertext
	}

	return pkcs5Validate(msg[aes.BlockSize:])
}

func aesDecryptCTR(key, msg []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := msg[:aes.BlockSize]
	msg = msg[aes.BlockSize:]

	counter := big.NewInt(0).SetBytes(iv)
	for i := 0; i < len(msg); i += aes.BlockSize {
		b := counter.Bytes()
		cipher.Encrypt(b, b)
		copy(msg[i:i+aes.BlockSize], xor(msg[i:i+aes.BlockSize], b))
		counter = counter.Add(counter, big.NewInt(1))
	}

	return msg[:], nil
}

func xor(a, b []byte) []byte {
	var result []byte
	if len(a) < len(b) {
		result = make([]byte, len(a))
		for i := range a {
			result[i] = a[i] ^ b[i]
		}
	} else {
		result = make([]byte, len(b))
		for i := range b {
			result[i] = a[i] ^ b[i]
		}
	}

	return result
}

func mustBytes(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}

	return b
}

func main() {
	keys := []string{
		"140b41b22a29beb4061bda66b6747e14",
		"140b41b22a29beb4061bda66b6747e14",
	}
	ciphertexts := []string{
		"4ca00ff4c898d61e1edbf1800618fb2828a226d160dad07883d04e008a7897ee2e4b7465d5290d0c0e6c6822236e1daafb94ffe0c5da05d9476be028ad7c1d81",
		"5b68629feb8606f9a6667670b75b38a5b4832d0f26e1ab7da33249de7d4afc48e713ac646ace36e872ad5fb8a512428a6e21364b0c374df45503473c5242a253",
	}

	for i := 0; i < len(keys); i++ {
		key := mustBytes(hex.DecodeString(keys[i]))
		msg := mustBytes(hex.DecodeString(ciphertexts[i]))
		d, err := aesDecryptCBC(key, msg)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(d))
	}

	keys = []string{
		"36f18357be4dbd77f050515c73fcf9f2",
		"36f18357be4dbd77f050515c73fcf9f2",
	}
	ciphertexts = []string{
		"69dda8455c7dd4254bf353b773304eec0ec7702330098ce7f7520d1cbbb20fc388d1b0adb5054dbd7370849dbf0b88d393f252e764f1f5f7ad97ef79d59ce29f5f51eeca32eabedd9afa9329",
		"770b80259ec33beb2561358a9f2dc617e46218c0a53cbeca695ae45faa8952aa0e311bde9d4e01726d3184c34451",
	}

	for i := 0; i < len(keys); i++ {
		key := mustBytes(hex.DecodeString(keys[i]))
		msg := mustBytes(hex.DecodeString(ciphertexts[i]))
		d, err := aesDecryptCTR(key, msg)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(d))
	}
}
