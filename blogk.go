package blogk

import (
	"log"
	// "crypto/rsa"
	// "crypto/sha1"
	"crypto/sha256"
	// "crypto/x509"
	// "encoding/base64"
	"encoding/json"
	"fmt"
	// "strings"

	"github.com/erikreppel/blogk/types"
	"github.com/pkg/errors"
	"github.com/tendermint/go-db"
	"github.com/tendermint/go-merkle"
	tmspTypes "github.com/tendermint/tmsp/types"
)

// Constants for a switch
const (
	UserConst = "USER="
	PostConst = "POST="
)

// NewBlogK returns a new BlogK node
func NewBlogK() *Blogk {
	usersDB, err := db.NewLevelDB("usersDB")
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
// tx will be in the form COMMAND=PubKey{}
func (blogk *Blogk) AppendTx(tx []byte) tmspTypes.Result {
	command := string(tx)[:5]
	body := []byte(string(tx)[5:])
	log.Println("parts:", string(body))

	switch command {
	case UserConst:
		err := blogk.insertUser(body)
		if err != nil {
			return tmspTypes.ErrBaseInvalidInput
		}
		return tmspTypes.OK
	default:
		return tmspTypes.ErrBaseInsufficientGasPrice
	}

}

// CheckTx just always returns true for now
func (blogk *Blogk) CheckTx(tx []byte) tmspTypes.Result {
	return tmspTypes.OK
}

// Commit returns a hash of the sum of the users and posts tree hashes
func (blogk *Blogk) Commit() tmspTypes.Result {
	concat := string(blogk.users.Hash()) + string(blogk.posts.Hash())
	hasher := sha256.New()
	_, err := hasher.Write([]byte(concat))
	if err != nil {
		return tmspTypes.ErrEncodingError
	}
	hashed := hasher.Sum(nil)
	return tmspTypes.NewResultOK(hashed, "")
}

// Query verifies if something is in the blockchain in log(n) time
func (blogk *Blogk) Query(query []byte) tmspTypes.Result {
	index, value, exists := blogk.users.Get(query)
	resStr := fmt.Sprintf("Index=%v value=%v exists=%v", index, string(value), exists)
	return tmspTypes.NewResultOK([]byte(resStr), "")
}

func (blogk *Blogk) insertUser(raw []byte) error {
	user := &types.User{}
	err := json.Unmarshal(raw, user)
	if err != nil {
		return err
	}
	err = verifyIDHashMatch(user)
	if err != nil {
		return err
	}
	blogk.users.Set(user.ID, user.Key)
	return nil

}

// where k is the public key and sig is the signed thing
func verifySigning(k, sig []byte) error {

	return nil
}

func verifyIDHashMatch(user *types.User) error {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(user.Key))
	if err != nil {
		return errors.Wrap(err, "Failed to write to hasher")
	}

	if string(hasher.Sum(nil)) != string(user.ID) {
		return errors.New("The hash of the key did not match the user ID")
	}
	return nil
}
