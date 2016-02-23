package pbkdf2

/*
	This package written by Kyle Isom was pulled from Github, so I've included it
    in my local source code. It's just a convenience wrapper around the Go PBKDF2 package.
*/

import (
	_pbkdf2 "golang.org/x/crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"fmt"
)

var (
	HashFunc   = sha1.New // Hash function to use
	Iterations = 16384    // Number of iterations to use
	KeySize    = 32       // Output size of key (32 is suitable for AES256)
	SaltSize   = 16       // Size of salt. Recommend >8.
)

// Type PasswordHash stores a hashed version of a password.
type PasswordHash struct {
	Hash []byte
	Salt []byte
}

func generateSalt(chars int) (salt []byte) {
	saltBytes := make([]byte, chars)
	nRead, err := rand.Read(saltBytes)
	if err != nil {
		salt = []byte{}
	} else if nRead < chars {
		salt = []byte{}
	} else {
		salt = saltBytes
	}
	return
}

// HashPassword generates a salt and returns a hashed version of the password.
func HashPassword(password string) *PasswordHash {
	salt := generateSalt(SaltSize)
	return HashPasswordWithSalt(password, salt)
}

// HashPasswordWithSalt hashes the password with the specified salt.
func HashPasswordWithSalt(password string, salt []byte) (ph *PasswordHash) {
	hash := _pbkdf2.Key([]byte(password), salt, Iterations, KeySize, HashFunc)
	return &PasswordHash{hash, salt}
}

// MatchPassword compares the input password with the password hash.
// It returns true if they match.
func MatchPassword(password string, ph *PasswordHash) bool {
	matched := 0
	new_hash := HashPasswordWithSalt(password, ph.Salt)
	fmt.Println("new hash is", new_hash)

	size := len(new_hash.Hash)
	if size > len(ph.Hash) {
		size = len(ph.Hash)
	}

	for i := 0; i < size; i++ {
		matched += subtle.ConstantTimeByteEq(new_hash.Hash[i], ph.Hash[i])
	}

	passed := matched == size
	if len(new_hash.Hash) != len(ph.Hash) {
		return false
	}
	return passed
}
