// Package password provides function to work with password and password hashes.
//
// It contains crypto logic in one place and provide convenient interface for
// other parts of the app. Main functions are `HashAndEncode` for generating hash
// from a plaintext password and `Check` that compares plaintext password
// against encoded string created with `HashAndEncode`.
//
// This packages uses Argon2 for password hashing with unique per-password salt
// and hash params that are described in package constants.
package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const (
	encodingHashType = "argon2id"

	// Default parameters for hashing
	DefaultIterations = 5
	DefaultMemory     = 32 * 1024 // 32MiB
	DefaultThreads    = 1
	DefaultLen        = 64
)

// genSalt generates random salt of n bytes size
func genSalt(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// HashParams describe parameters for Argon2 hashing.
// These parameters are encoded in the Encode function.
type HashParams struct {
	Salt       []byte
	Iterations uint32
	Memory     uint32
	Threads    uint8
	Len        uint32
}

// NewHashParams creates hash params struct with default values
func NewHashParams() (HashParams, error) {
	salt, err := genSalt(16)
	if err != nil {
		return HashParams{}, err
	}

	return HashParams{
		Salt:       salt,
		Iterations: DefaultIterations,
		Memory:     DefaultMemory,
		Threads:    DefaultThreads,
		Len:        DefaultLen,
	}, nil
}

// Encode converts hash from byte slice to the string with hash params
// including salt
func Encode(hash []byte, params HashParams) string {
	base64Salt := base64.RawStdEncoding.EncodeToString(params.Salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Threads,
		base64Salt,
		base64Hash,
	)
}

// Decode converts hash from encoded string to the byte slice of hashed password
// and the struct of hash params
func Decode(s string) ([]byte, HashParams, error) {
	// $argon2id$v=<version>$m=<memory>,t=<iters>,p=<threads>$<salt>$<hash>
	vals := strings.Split(s, "$")

	if vals[1] != encodingHashType {
		return nil, HashParams{}, fmt.Errorf("invalid hash type, expected %s, got %s\n", encodingHashType, vals[1])
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, HashParams{}, errors.Wrapf(err, "failed to parse password hash version")
	}

	if version != argon2.Version {
		return nil, HashParams{}, fmt.Errorf("invalid password hash version %d", version)
	}

	var params HashParams

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Threads)
	if err != nil {
		return nil, HashParams{}, errors.Wrapf(err, "failed to parse password hash params")
	}

	params.Salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, HashParams{}, errors.Wrapf(err, "failed to base64 decode password hash salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, HashParams{}, errors.Wrapf(err, "failed to base64 decode password hash")
	}

	params.Len = uint32(len(hash))

	return hash, params, nil
}

// genHash generates password hash for a given plaintext password
// using given params
func HashWithParams(password string, params HashParams) []byte {
	return argon2.IDKey(
		[]byte(password),
		params.Salt,
		params.Iterations,
		params.Memory,
		params.Threads,
		params.Len,
	)
}

// HashEncode creates Argon2 password hash for a given plaintext password
// using default hash parameters
func Hash(password string) ([]byte, HashParams, error) {
	params, err := NewHashParams()
	if err != nil {
		return nil, HashParams{}, errors.Wrap(err, "failed to create password hash params")
	}

	hash := HashWithParams(password, params)

	return hash, params, nil
}

// HashEncode creates Argon2 password hash for a given plaintext password
// using default hash parameters and returns its encoded representation
func HashAndEncode(password string) (string, error) {
	hash, params, err := Hash(password)
	if err != nil {
		return "", errors.Wrap(err, "failed to create password hash")
	}

	return Encode(hash, params), nil
}

// Check takes plaintext password and compares it to the encoded hash value
func Check(password string, encoded string) (bool, error) {
	decodedHash, decodedParams, err := Decode(encoded)
	if err != nil {
		return false, errors.Wrap(err, "failed to decode password hash for compare")
	}

	hash := HashWithParams(password, decodedParams)

	if subtle.ConstantTimeCompare(hash, decodedHash) == 1 {
		return true, nil
	}

	return false, nil
}
