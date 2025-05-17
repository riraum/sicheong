package security

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := &[32]byte{}

	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		plaintxt []byte
		key      *[32]byte
	}{
		{
			plaintxt: []byte("Test txt!"),
			key:      key,
		},
	}

	for _, test := range tests {
		ciphertxt, err := Encrypt(test.plaintxt, test.key)
		if err != nil {
			t.Fatal(err)
		}

		plaintxt, err := Decrypt(ciphertxt, test.key)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(plaintxt, test.plaintxt) {
			t.Errorf("plaintxts don't match")
		}
		// ciphertxt[0] ^= 0xdd
		// plaintxt, err = Decrypt(ciphertxt, test.key)
		//
		//	if err == nil {
		//		t.Errorf("gcmOpen should not have worked, but did")
		//	}
	}
}
