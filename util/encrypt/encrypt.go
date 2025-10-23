package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var (
	keyOnce sync.Once
	key     []byte
	keyErr  error
)

// SecureString transparently encrypts/decrypts PHI at BSON boundary (AES-256-GCM).
// Stored format: v1:gcm:<base64(nonce|ciphertext)>
type SecureString string

func loadKey() error {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		return errors.New("ENCRYPTION_KEY not set")
	}
	k, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return err
	}
	if len(k) != 32 {
		return errors.New("ENCRYPTION_KEY must decode to 32 bytes")
	}
	key = k
	return nil
}

func ensureKey() error {
	keyOnce.Do(func() {
		keyErr = loadKey()
	})
	return keyErr
}

func (s SecureString) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if err := ensureKey(); err != nil {
		return bsontype.String, nil, err
	}
	plain := string(s)
	if plain == "" {
		return bson.MarshalValue("")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return 0, nil, err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return 0, nil, err
	}
	ct := aead.Seal(nil, nonce, []byte(plain), nil)
	payload := "v1:gcm:" + base64.StdEncoding.EncodeToString(append(nonce, ct...))
	return bson.MarshalValue(payload)
}

func (s *SecureString) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var enc string
	if err := bson.UnmarshalValue(t, data, &enc); err != nil {
		return err
	}
	if enc == "" {
		*s = ""
		return nil
	}
	if !strings.HasPrefix(enc, "v1:gcm:") {
		// fallback: assume plain string stored
		*s = SecureString(enc)
		return nil
	}
	if err := ensureKey(); err != nil {
		return err
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(enc, "v1:gcm:"))
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	if len(raw) < aead.NonceSize() {
		return errors.New("ciphertext too short")
	}
	nonce := raw[:aead.NonceSize()]
	ct := raw[aead.NonceSize():]
	pt, err := aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return err
	}
	*s = SecureString(string(pt))
	return nil
}
