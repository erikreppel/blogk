package utils

import (
	"crypto/ecdsa"
	// "crypto/elliptic"
	// "crypto/md5"
	// "crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"sync"
	// "hash"
	// "io"
	"encoding/json"
	// "math/big"

	"github.com/erikreppel/blogk"
	"github.com/erikreppel/blogk/types"
	"github.com/tendermint/tmsp/client"
)

// AddUser adds a user to the local cluster
func AddUser(pubkey *ecdsa.PublicKey) {
	fmt.Println("Public Key :")
	fmt.Printf("%x \n", pubkey)

	pub, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()
	_, err = hasher.Write(pub)
	if err != nil {
		panic(err)
	}

	user := &types.User{
		Key: pub,
		ID:  hasher.Sum(nil),
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	payload := []byte("USER=")
	payload = append(payload, userJSON...)
	mut := &sync.Mutex{}
	blk := blogk.NewBlogK()
	client := tmspcli.NewLocalClient(mut, blk)
	fmt.Println("Info:", client.Info())
	fmt.Println(client.AppendTxSync(payload))
	client.Stop()
	fmt.Println("Client running:", client.IsRunning())
}
