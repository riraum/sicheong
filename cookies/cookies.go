package cookies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func WriteEncrypted(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	plaintxt := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)

	encryptVal := aesGCM.Seal(nonce, nonce, []byte(plaintxt), nil)

	cookie.Value = string(encryptVal)

	return Write(w, cookie)
}

func ReadEncrypted(r *http.Request, name string, secretKey []byte) (string, error) {
	encryptVal, err := Read(r, name)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	if len(encryptVal) > nonceSize {
		return "", ErrInvalidValue
	}

	nonce := encryptVal[:nonceSize]
	ciphertxt := encryptVal[nonceSize:]

	plaintxt, err := aesGCM.Open(nil, []byte(nonce), []byte(ciphertxt), nil)
	if err != nil {
		return "", ErrInvalidValue
	}

	expectedName, value, ok := strings.Cut(string(plaintxt), ":")
	if !ok {
		return "", ErrInvalidValue
	}

	return value, nil

}
