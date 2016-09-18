package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/go-db"
	"github.com/tendermint/go-merkle"
	"github.com/tendermint/tmsp/types"
)

const (
	User = "USER"
	Post = "POST"
)

// NewBlogK returns a new BlogK node
func NewBlogK() *Blogk {
	usersDB, err := db.NewLevelDB("userDB")
	if err != nil {
		panic(err)
	}
	postsDB, err := db.NewLevelDB("postsDB")
	if err != nil {
		panic(err)
	}
	return &Blogk{
		users:   merkle.NewIAVLTree(100000, usersDB),
		posts:   merkle.NewIAVLTree(100000, postsDB),
		options: make(map[string]string),
	}

}

// Blogk is the base structure of the validator node
type Blogk struct {
	users   merkle.Tree
	posts   merkle.Tree
	options map[string]string
}

// Info returns info about the users and post contained in the blogk
func (blogk *Blogk) Info() string {
	return fmt.Sprintf("Number of users: %d, Number of posts: %d", blogk.users.Size(), blogk.posts.Size())
}

// SetOption sets configs in the blogk
func (blogk *Blogk) SetOption(key, value string) (log string) {
	blogk.options[key] = value
	return fmt.Sprintf("Current options: %+v", blogk.options)
}

// AppendTx adds a post of the chain after verifying it was sent by a proper key
func (blogk *Blogk) AppendTx(tx []byte) (types.Result, []byte) {
	parts := strings.Split(string(tx), "=")
	if len(parts) != 2 {
		return types.ErrBaseInvalidInput, nil
	}
	switch parts[0] {
	case User:
		id, err := blogk.insertUser([]byte(parts[1]))
		if err != nil {
			return types.ErrBaseInvalidInput, nil
		}
		return types.OK, id
	case Post:

	}

}

func (blogk *Blogk) insertUser(raw []byte) (error, []byte) {
	user := &User{}
	err := json.Unmarshal(raw, user)
	if err != nil {
		return err, nil
	}
	err = verifyIDHashMatch(user)
	if err != nil {
		return err, nil
	}

}

func verifySigning(k, sig []byte) error {
	key, err := base64.StdEncoding.DecodeString(string(k))
	if err != nil {
		return err
	}
	re, err := x509.ParsePKIXPublicKey(key)
	pub := re.(*rsa.PublicKey)
	if err != nil {
		return err
	}
	// check signing

}

func verifyIDHashMatch(user *User) error {
	hasher := sha256.New()
	_, err = hasher.Write([]byte(user.Key))
	if err != nil {
		return errors.Wrap(err, "Failed to write to hasher")
	}
	if hasher.Sum() != user.ID {
		return errors.New("The hash of the key did not match the user ID")
	}
	return nil
}
