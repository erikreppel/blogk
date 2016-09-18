package types

// User is what gets sent to register a new user.
// Contains Key wich is the public key as bytes, and ID, a hash of the public key
type User struct {
	Key []byte `json:"key"`
	ID  []byte `json:"id"`
}

// Post represents a message to be inserted into the chain
type Post struct {
	R        int
	SignHash string
}
