package types

import (
	"crypto/ecdsa"
	"crypto/x509"
	"math/big"

	"github.com/pkg/errors"
)

// User is what gets sent to register a new user.
// Contains Key wich is the public key as bytes, and ID, a hash of the public key
type User struct {
	Key []byte `json:"key"`
	ID  []byte `json:"id"`
}

// Post represents a message to be inserted into the chain
type Post struct {
	R         big.Int
	S         big.Int
	Signed    []byte
	Raw       []byte
	UserID    []byte
	Signature []byte
}

// Verify ensures Signed is a raw signed by a private key
func (p *Post) Verify(key []byte) error {
	pk, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return errors.Wrap(err, "Failed to load public key")
	}
	pubkey := pk.(*ecdsa.PublicKey)
	if ecdsa.Verify(pubkey, p.Signed, &p.R, &p.S) {
		return nil
	}
	return errors.New("Failed to verify that the correct private key signed the message")
}
