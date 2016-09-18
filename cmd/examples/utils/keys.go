package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	// "crypto/md5"
	"crypto/rand"
	// "crypto/sha256"
	// "crypto/x509"
	// "fmt"
	// "sync"
	// "hash"
	// "io"
	// "encoding/json"
	// "math/big"
)

// MakeKeys returns public and private keys to be used
func MakeKeys() (*ecdsa.PublicKey, *ecdsa.PrivateKey) {
	pubkeyCurve := elliptic.P256()
	privatekey := new(ecdsa.PrivateKey)
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey := privatekey.Public().(*ecdsa.PublicKey)
	return pubkey, privatekey
}
