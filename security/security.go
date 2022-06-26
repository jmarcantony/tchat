package security

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/box"
)

func GenerateKey() ([]byte, []byte) {
	publicKey, privateKey, _ := box.GenerateKey(rand.Reader)
	return (*publicKey)[:], (*privateKey)[:]
}
