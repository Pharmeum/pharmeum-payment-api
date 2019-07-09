package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

func newAddress() (string, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", err
	}

	r, s, err := ecdsa.Sign(rand.Reader, private, md5.New().Sum(nil))
	if err != nil {
		return "", err
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	return fmt.Sprintf("%x", signature), nil
}
