package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
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
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap cipher block: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to read nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintxt, nil), nil
}

func Decrypt(ciphertxt []byte, key *[32]byte) (plaintxt []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap cipher block: %w", err)
	}

	if len(ciphertxt) < gcm.NonceSize() {
		return nil, fmt.Errorf("mailformed ciphertxt: %w", err)
	}

	return gcm.Open(nil,
		ciphertxt[:gcm.NonceSize()],
		ciphertxt[gcm.NonceSize():],
		nil,
	)
}

func Hash(tag string, data []byte) []byte {
	h := hmac.New(sha512.New512_256, []byte(tag))
	h.Write(data)

	return h.Sum(nil)
}

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 14) //nolint:mnd
}

func CheckPasswordHash(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
