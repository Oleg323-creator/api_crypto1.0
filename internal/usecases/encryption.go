package usecases

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

// ENCRYPTIONS USING AESGCM
func (u *Usecases) EncryptAESGCM(password, plaintext []byte) (string, string, error) {
	block, err := aes.NewCipher(password)
	if err != nil {
		return "", "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	log.Printf("Generated Nonce : %x", nonce)

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	base64Nonce := hex.EncodeToString(nonce)
	base64ciphertext := hex.EncodeToString(ciphertext)

	log.Printf("Generated Ciphertext: %x", ciphertext)

	return base64ciphertext, base64Nonce, nil
}

func (u *Usecases) DecryptAESGCM(nonceBase64 string, base64ciphertext string, password []byte) ([]byte, error) {

	nonce, err := hex.DecodeString(nonceBase64)
	if err != nil {
		log.Fatalf("Error decoding nonce: %v", err)
	}

	ciphertext, err := hex.DecodeString(base64ciphertext)
	if err != nil {
		log.Fatalf("Error decoding ciphertext: %v", err)
	}

	block, err := aes.NewCipher(password)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Printf("Decryption failed: %v", err)
		return nil, err
	}

	log.Printf("Decrypted text: %x", string(plaintext))

	return plaintext, nil
}
