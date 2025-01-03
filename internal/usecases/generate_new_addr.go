package usecases

import (
	"api_crypto1.0/internal/db/repository"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func (u *Usecases) GenerateNewAdd(curr string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicAddres := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	fmt.Printf("Public address: %s", publicAddres)

	derKey := crypto.FromECDSA(privateKey)

	password := []byte(os.Getenv("SECRET_PASSWORD"))
	if os.Getenv("SECRET_PASSWORD") == "" {
		log.Fatal("Error getting secret password from .env")
	}
	log.Printf("Using secret password: %s", os.Getenv("SECRET_PASSWORD"))

	encryptedKey, nonce, err := u.EncryptAESGCM(password, derKey)
	if err != nil {
		log.Fatal(err)
	}

	data := repository.DataToSave{
		PrivateKey: encryptedKey,
		Address:    publicAddres,
		Currency:   curr,
		Nonce:      nonce,
	}

	err = u.Repository.SaveNewAddr(data)
	if err != nil {
		return "", fmt.Errorf("Error saiving new addr to DB: %v", err)
	}

	return publicAddres, nil
}
