package battletoken

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

// BattleToken は AES-128-ECB で monster UUID を暗号化/復号する。
// UUID は 16 バイトなのでちょうど 1 ブロック。
type BattleToken struct {
	key [16]byte
}

func New(secret string) *BattleToken {
	hash := sha256.Sum256([]byte(secret))
	var key [16]byte
	copy(key[:], hash[:16])
	return &BattleToken{key: key}
}

func (t *BattleToken) Encrypt(id uuid.UUID) (string, error) {
	block, err := aes.NewCipher(t.key[:])
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, 16)
	block.Encrypt(ciphertext, id[:])
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(ciphertext), nil
}

func (t *BattleToken) Decrypt(token string) (uuid.UUID, error) {
	ciphertext, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(token)
	if err != nil || len(ciphertext) != 16 {
		return uuid.Nil, fmt.Errorf("invalid battle token")
	}
	block, err := aes.NewCipher(t.key[:])
	if err != nil {
		return uuid.Nil, err
	}
	plaintext := make([]byte, 16)
	block.Decrypt(plaintext, ciphertext)
	var id uuid.UUID
	copy(id[:], plaintext)
	return id, nil
}
