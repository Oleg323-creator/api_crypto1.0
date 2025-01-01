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

// EncryptAESGCM шифрует данные с использованием AES-GCM и возвращает ciphertext и строковое представление nonce
func (u *Usecases) EncryptAESGCM(password, plaintext []byte) (string, string, error) {
	block, err := aes.NewCipher(password)
	if err != nil {
		return "", "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// Генерация случайного nonce
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	log.Printf("Generated Nonce : %x", nonce)

	// Шифрование данных
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// Преобразуем nonce в строку Base64
	base64Nonce := hex.EncodeToString(nonce)
	base64ciphertext := hex.EncodeToString(ciphertext)

	log.Printf("Generated Ciphertext: %x", ciphertext)

	return base64ciphertext, base64Nonce, nil
}

// DecryptAESGCM расшифровывает данные с использованием AES-GCM
func (u *Usecases) DecryptAESGCM(nonceBase64 string, base64ciphertext string, password []byte) ([]byte, error) {
	// Декодируем строку Base64 обратно в байтовый массив для nonce
	nonce, err := hex.DecodeString(nonceBase64)
	if err != nil {
		log.Fatalf("Ошибка при декодировании Base64: %v", err)
	}

	ciphertext, err := hex.DecodeString(base64ciphertext)
	if err != nil {
		log.Fatalf("Ошибка при декодировании Base64: %v", err)
	}

	log.Printf("Decoded Ciphertext =: %x", ciphertext)

	log.Printf("Decoded Nonce =: %x", nonce)

	// Создаем шифратор
	block, err := aes.NewCipher(password)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Создаем GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Расшифровка
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Printf("Decryption failed: %v", err)
		return nil, err
	}

	// Логируем расшифрованный текст
	log.Printf("Decrypted text: %x", string(plaintext))

	return plaintext, nil
}
