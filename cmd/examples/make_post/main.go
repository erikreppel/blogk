package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	// "math/big"
	"fmt"
	"sync"
	"time"

	"github.com/erikreppel/blogk"
	"github.com/erikreppel/blogk/cmd/examples/utils"
	"github.com/erikreppel/blogk/types"
	"github.com/tendermint/tmsp/client"
)

type userPost struct {
	Username  string `json:"username"`
	UserID    string `json:"user_id"`
	Body      string `json:"body"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	mut := &sync.Mutex{}
	blk := blogk.NewBlogK()
	client := tmspcli.NewLocalClient(mut, blk)

	pub, priv := utils.MakeKeys()
	fmt.Println("Make keys")
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		panic(err)
	}
	hasher := sha256.New()
	_, err = hasher.Write(pubBytes)
	if err != nil {
		panic(err)
	}
	fmt.Println("Public Key :")
	fmt.Printf("%x \n", pub)

	// Create user
	fmt.Println("Creating user")
	user := &types.User{
		Key: pubBytes,
		ID:  hasher.Sum(nil),
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	payload := []byte("USER=")
	payload = append(payload, userJSON...)
	fmt.Println("Info:", client.Info())
	fmt.Println(client.AppendTxSync(payload))
	fmt.Println("Client running:", client.IsRunning())

	// Make post
	fmt.Println("Making post")
	userID := hasher.Sum(nil)
	up := &userPost{
		Username:  "Erik",
		UserID:    string(hasher.Sum(nil)),
		Timestamp: time.Now().Unix(),
		Body:      "Hello, world!",
	}
	m, _ := json.Marshal(up)
	r, s, serr := ecdsa.Sign(rand.Reader, priv, m)
	if serr != nil {
		panic(err)
	}
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	hasher = sha256.New()
	_, err = hasher.Write(m)
	if err != nil {
		panic(err)
	}
	post := &types.Post{
		R:         *r,
		S:         *s,
		Raw:       m,
		UserID:    userID,
		Hashed:    hasher.Sum(nil),
		Signature: signature,
	}

	mp, _ := json.Marshal(post)
	tmpP := string(userID) + "." + string(post.Hashed) + "." + string(mp)
	payload2 := append([]byte("POST="), tmpP...)
	fmt.Println(client.AppendTxSync(payload2))

}
