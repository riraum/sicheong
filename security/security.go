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

func NewEncryptionKey() (*[32]byte, error) {
	key := [32]byte{}

	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		return nil, fmt.Errorf("failed to generate random encryption key: %w", err)
	}

	return &key, nil
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

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
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

	decrypted, err := gcm.Open(nil,
		ciphertxt[:gcm.NonceSize()],
		ciphertxt[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return decrypted, fmt.Errorf("failed to gcmOpen: %w", err)
	}

	return decrypted, nil
}

func Hash(tag string, data []byte) []byte {
	h := hmac.New(sha512.New512_256, []byte(tag))
	h.Write(data)

	return h.Sum(nil)
}

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 14) //nolint:mnd
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return hashedPassword, nil
}

func CheckPasswordHash(hash, password []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return fmt.Errorf("failed to check hashed password: %w", err)
	}

	return nil
}
