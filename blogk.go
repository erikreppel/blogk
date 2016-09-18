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
	"strings"

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
		userID, err := blogk.insertUser(body)
		if err != nil {
			log.Println("ERROR:", err)
			return tmspTypes.ErrBaseInvalidInput
		}
		l := fmt.Sprintln("Successfully inserted as id", string(userID))
		return tmspTypes.NewResultOK(userID, l)
	case PostConst:
		id, err := blogk.insertPost(body)
		if err != nil {
			log.Println("ERROR:", err)
			return tmspTypes.ErrBaseInvalidInput
		}
		l := fmt.Sprintln("Successfully inserted post as id", string(id))
		return tmspTypes.NewResultOK(id, l)
	default:
		err := fmt.Errorf("unsupported %s", command)
		log.Println("ERROR: ", err)
		return tmspTypes.ErrBaseInvalidInput
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
	if !exists {
		return tmspTypes.ErrBadNonce
	}
	resStr := fmt.Sprintf("Index=%v value=%v exists=%v", index, string(value), exists)
	log.Println(resStr)
	return tmspTypes.NewResultOK([]byte(value), "")
}

func (blogk *Blogk) insertPost(raw []byte) ([]byte, error) {
	splt := strings.Split(string(raw), ".")
	userID := splt[0]
	postHash := splt[1]
	rawPost := splt[2]
	// get public key for that user
	resp := blogk.Query([]byte(userID))
	if resp.Code.String() == tmspTypes.ErrBadNonce.Code.String() {
		log.Println("ERROR: bad nonce")
		return nil, errors.New("Public Key not found")
	} else if len(resp.Data) == 0 {
		log.Println("ERROR: Got no data back when looking up public key")
		return nil, errors.New("Public Key not found")
	}
	key := resp.Data
	post := &types.Post{}
	err := json.Unmarshal([]byte(rawPost), post)
	if err != nil {
		e := errors.Wrap(err, "Failed to Unmarshal json")
		log.Println(e)
		return nil, e
	}
	if userID != string(post.UserID) {
		err := errors.New("UserID from raw post does not match that in the body")
		return nil, err
	}
	// verify post origins
	err = post.Verify(key)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to verify poster")
	}
	postID := []byte(string(post.UserID) + "." + postHash)
	blogk.posts.Set(postID, post.Raw)
	return postID, nil
}

func (blogk *Blogk) insertUser(raw []byte) ([]byte, error) {
	user := &types.User{}
	err := json.Unmarshal(raw, user)
	if err != nil {
		return nil, err
	}
	err = verifyIDHashMatch(user)
	if err != nil {
		return nil, err
	}
	blogk.users.Set(user.ID, user.Key)
	return user.ID, nil
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
