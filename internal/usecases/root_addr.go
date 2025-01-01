package usecases

import (
	"api_crypto1.0/internal/db/repository"
	"context"
	"crypto/x509"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
)

func (u *Usecases) CheckRootAddrInDB() (string, error) {
	// CHECKING IF EXISTS

	addr, err := u.Repository.GetRootAddr()
	if err != nil {
		//	log.Printf("No root address in db, creating new addr: %v", err)
		addr, err = u.GenerateRootAddr()
		if err != nil {
			return "", fmt.Errorf("Error generating new root address: %v", err)
		}
	}

	log.Printf("Root address: %s", addr)

	return addr, nil
}

func (u *Usecases) GenerateRootAddr() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("Error loading .env file: %v", err)
	}

	privateKey1, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKey := crypto.FromECDSA(privateKey1)

	log.Printf("PRIVATE KEY1: %x", privateKey)

	publicAddres := crypto.PubkeyToAddress(privateKey1.PublicKey).Hex()
	log.Printf("Public address: %s", publicAddres)

	// CONVERT KEY TO DER FORMAT
	derKey := crypto.FromECDSA(privateKey1)

	password := []byte(os.Getenv("SECRET_PASSWORD"))
	if len(password) == 0 {
		return "", fmt.Errorf("SECRET_PASSWORD is not set in .env")
	}

	encryptedKey, nonce, err := u.EncryptAESGCM(password, derKey)
	if err != nil {
		log.Fatalf("Error encrypting key: %v", err)
	}

	data := repository.DataToSave{
		PrivateKey: encryptedKey,
		Address:    publicAddres,
		Nonce:      nonce,
	}

	err = u.Repository.SaveRootAddrToDB(data)
	if err != nil {
		return "", fmt.Errorf("Error saiving new addr to DB: %v", err)
	}

	return publicAddres, nil
}

func (u *Usecases) MergeCoinsToRoot(data repository.Params) error {

	fromAddress := data.ToAddr
	toAddress, err := u.Repository.GetRootAddr()
	if err != nil {
		return err
	}
	amount := data.Value

	bigIntAmount := new(big.Int)

	if _, success := bigIntAmount.SetString(amount, 10); success {
		fmt.Printf("BigInt value: %d", bigIntAmount)
	} else {
		fmt.Println("Failed to convert string to *big.Int")
	}

	gasPrice, err := u.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Error getting gas price: %v", err)
	}

	nonce, err := u.Client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		log.Fatalf("Error getting nonce: %v", err)
	}

	gasLimit := uint64(21000)

	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), bigIntAmount, gasLimit, gasPrice, nil)

	password := []byte(os.Getenv("SECRET_PASSWORD"))
	if len(password) != 32 {
		log.Fatalf("Password length should be 32 bytes. Got: %d", len(password))
	}

	log.Printf("Using secret password: %s", os.Getenv("SECRET_PASSWORD"))

	privateKeyToDecription, nonc, err := u.Repository.GetPrivateKeyFromDB(fromAddress)

	privateKey, err := u.DecryptAESGCM(nonc, privateKeyToDecription, password)
	if err != nil {
		return fmt.Errorf("Errore decription private key: %v", err)
	}

	privateKeyInCorrType, err := x509.ParseECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP2930Signer(big.NewInt(11155111)), privateKeyInCorrType)
	if err != nil {
		log.Fatalf("Error singing tx: %v", err)
	}

	err = u.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Error sending tx: %v", err)
	}

	log.Printf("Tx has sent! Hash: %s\n", signedTx.Hash().Hex())
	log.Println("//////////////////////////////////////////////////////////////////////////////")

	return nil
}
