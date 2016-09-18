package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	// "crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"sync"
	// "hash"
	// "io"
	"encoding/json"
	"math/big"

	"github.com/erikreppel/blogk"
	"github.com/erikreppel/blogk/types"
	"github.com/tendermint/tmsp/client"
)

func main() {
	pubkeyCurve := elliptic.P256()
	privatekey := new(ecdsa.PrivateKey)
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey := privatekey.Public().(*ecdsa.PublicKey)

	fmt.Println("Private Key :")
	fmt.Printf("%x \n", privatekey)
	r := big.NewInt(0)
	s := big.NewInt(0)
	fmt.Println("Public Key :")
	fmt.Printf("%x \n", pubkey)

	testStr := []byte("This is a test")

	r, s, serr := ecdsa.Sign(rand.Reader, privatekey, testStr)
	if serr != nil {
		panic(err)
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	verifystatus := ecdsa.Verify(pubkey, testStr, r, s)
	fmt.Println("Verified:", verifystatus)
	pub, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(pub))

	k, err := x509.ParsePKIXPublicKey(pub)
	if err != nil {
		panic(err)
	}
	pubkey2 := k.(*ecdsa.PublicKey)
	verifystatus = ecdsa.Verify(pubkey2, testStr, r, s)
	fmt.Println("Verified:", verifystatus)

	hasher := sha256.New()
	_, err = hasher.Write(pub)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pub))

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
	fmt.Println("PAYLOAD:", payload)
	mut := &sync.Mutex{}
	blk := blogk.NewBlogK()
	client := tmspcli.NewLocalClient(mut, blk)
	fmt.Println("Info:", client.Info())
	fmt.Println(client.AppendTxSync(payload))

}
