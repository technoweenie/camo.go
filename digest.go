package camo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

type DigestCalculator interface {
	Calculate(value string) string
}

type Digest struct {
	key string
}

func NewDigest(key string) *Digest {
	return &Digest{key}
}

func (digest *Digest) Calculate(value string) string {
	hmac := hmac.New(sha1.New, []byte(digest.key))
	hmac.Write([]byte(value))
	return hex.EncodeToString([]byte(string(hmac.Sum(nil))))
}
