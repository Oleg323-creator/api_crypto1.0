package usecases

import (
	"api_crypto1.0/internal/db/repository"
	"context"
	"crypto/x509"
	"encoding/base64"
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
		log.Fatalf("Error loading .env file: %v", err)
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicAddres := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	log.Printf("Public address: %s", publicAddres)

	// CONVERT KEY TO DER FORMAT
	derKey := crypto.FromECDSA(privateKey)

	password := []byte(os.Getenv("SECRET_PASSWORD"))
	if os.Getenv("SECRET_PASSWORD") == "" {
		log.Fatal("Error getting secret password from .env")
	}

	encryptedKey, err := u.EncryptAESGCM(password, derKey)
	if err != nil {
		log.Fatalf("Error encrypting key: %v", err)
	}

	encryptedKeyBase64 := base64.StdEncoding.EncodeToString(encryptedKey)

	data := repository.DataToSave{
		PrivateKey: encryptedKeyBase64,
		Address:    publicAddres,
	}

	err = u.Repository.SaveRootAddrToDB(data)
	if err != nil {
		return "", fmt.Errorf("Error saiving new addr to DB: %v", err)
	}

	return publicAddres, nil
}

func (u *Usecases) MergeCoinsToRoot(data repository.TxData) error {

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
	if os.Getenv("SECRET_PASSWORD") == "" {
		log.Fatal("Error getting secret password from .env")
	}
	log.Printf("Using secret password: %s", os.Getenv("SECRET_PASSWORD"))

	privateKeyToDecription, err := u.Repository.GetPrivateKeyFromDB(fromAddress)

	privateKey, err := u.DecryptAESGCM(password, []byte(privateKeyToDecription))
	if err != nil {
		return fmt.Errorf("Errore decription private key: %v", err)
	}

	privateKeyInCorrType, err := x509.ParseECPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
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

	return nil
}
