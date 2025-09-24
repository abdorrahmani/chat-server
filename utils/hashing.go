package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
	"strings"
)

func HashMessage(message, algo, key string) (string, error) {
	var h hash.Hash

	switch strings.ToLower(algo) {
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	case "hmac-sha256":
		h = hmac.New(sha256.New, []byte(key))
	default:
		return "", errors.New("unknown hash algorithm")
	}

	h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil)), nil
}
