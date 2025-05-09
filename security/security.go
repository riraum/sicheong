package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

func Encrypt(plaintxt []byte, key *[32]byte) (ciphertxt []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintxt, nil), nil
}

func Decrypt(ciphertxt []byte, key *[32]byte) (plaintxt []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertxt) < gcm.NonceSize() {
		return nil, errors.New("mailformed ciphertxt")
	}
	return gcm.Open(nil,
		ciphertxt[:gcm.NonceSize()],
		ciphertxt[gcm.NonceSize():],
		nil,
	)
}
