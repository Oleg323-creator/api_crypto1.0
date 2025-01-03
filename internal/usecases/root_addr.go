package usecases

import (
	"api_crypto1.0/internal/db/repository"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
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
		log.Printf("No root address in db, creating new addr: %v", err)
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

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicAddres := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	log.Printf("Public address: %s", publicAddres)

	derKey := crypto.FromECDSA(privateKey)

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

	err = u.Repository.SaveRootData(data)
	if err != nil {
		return "", fmt.Errorf("Error saiving new addr to DB: %v", err)
	}

	return publicAddres, nil
}

func (u *Usecases) MergeCoinsToRoot(data repository.Params) error {
	log.Println("STARTING MERGING!")
	fromAddress := data.ToAddr
	toAddress, err := u.Repository.GetRootAddr()
	if err != nil {
		return err
	}

	balance, err := u.Client.BalanceAt(context.Background(), common.HexToAddress(fromAddress), nil)
	if err != nil {
		log.Fatalf("Error getting balance: %v", err)
	}

	toAddr := common.HexToAddress(toAddress)

	msg := ethereum.CallMsg{
		From: common.HexToAddress(fromAddress),
		To:   &toAddr,
		Gas:  0,
		Data: nil,
	}
	gasLimit, err := u.Client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("Error estimating gas: %v", err)
	}

	gasPrice, err := u.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Error getting gas price: %v", err)
	}

	gasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)

	// CHECKING IF ENOUGH BALANCE
	if balance.Cmp(gasCost) <= 0 {
		log.Fatalf("Insufficient funds")
	}

	amountToSend := new(big.Int).Sub(balance, gasCost)

	nonce, err := u.Client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		log.Fatalf("Error getting nonce: %v", err)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), amountToSend, gasLimit, gasPrice, nil)

	password := []byte(os.Getenv("SECRET_PASSWORD"))
	if len(password) != 32 {
		log.Fatalf("Password length should be 32 bytes. Got: %d", len(password))
	}

	privateKeyToDecription, nonc, err := u.Repository.GetPrivateKey(fromAddress)

	privateKey, err := u.DecryptAESGCM(nonc, privateKeyToDecription, password)
	if err != nil {
		return fmt.Errorf("Errore decription private key: %v", err)
	}

	privateKeyConverted, err := crypto.ToECDSA(privateKey)
	if err != nil {
		log.Fatalf("Error converting key: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP2930Signer(big.NewInt(11155111)), privateKeyConverted)
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
