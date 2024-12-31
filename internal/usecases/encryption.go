package usecases

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// ENCRIPTION FUNC USING AES-GCM
func (u *Usecases) EncryptAESGCM(password, plaintext []byte) ([]byte, error) {
	// CREATE AES KEY
	block, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	// VERY DIFFICULT ENCRIPTION PROCESS...
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// DECRYPTION
func (u *Usecases) DecryptAESGCM(password, ciphertext []byte) ([]byte, error) {

	block, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := ciphertext[:12], ciphertext[12:]

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
