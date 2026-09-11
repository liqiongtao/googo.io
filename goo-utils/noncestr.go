package goo_utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func NonceStr() string {
	bf := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, bf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bf)
}

func NonceStr8() string {
	bf := make([]byte, 4)
	if _, err := io.ReadFull(rand.Reader, bf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bf)
}
