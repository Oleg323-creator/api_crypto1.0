package usecases

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"os"
)

func (u *Usecases) Withdraw(addr string, amount string) error {

	fromAddress, err := u.Repository.GetRootAddr()
	if err != nil {
		return err
	}

	toAddress := addr

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
	log.Printf("Password length: %d", len(password))

	privateKeyToDecription, nonc, err := u.Repository.GetRootPrivateKey(fromAddress)

	privateKey, err := u.DecryptAESGCM(nonc, privateKeyToDecription, password)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key for address %s: %v", fromAddress, err)
	}
	log.Printf("PRIVATE KEY: %x", privateKey)
	/*
		privateKeyInCorrType, err := x509.ParseECPrivateKey(privateKey)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %v", err)
		}
	*/

	privateKey2, err := crypto.ToECDSA(privateKey)
	if err != nil {
		log.Fatalf("Ошибка при преобразовании ключа: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP2930Signer(big.NewInt(11155111)), privateKey2)
	if err != nil {
		log.Fatalf("Error singing tx: %v", err)
	}

	err = u.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Error sending tx: %v", err)
	}

	log.Printf("Money is comming to your address! Hash: %s\n", signedTx.Hash().Hex())
	log.Println("................................................................") //ШОБ В ЛОГАХ ВЫДЕЛЯЛОСЬ

	return nil
}

func (u *Usecases) GenerateWithdrawID() (int, error) {
	minId := 10000000
	maxId := 99999999

	rangeSize := maxId - minId + 1

	randomBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(rangeSize)))
	if err != nil {
		return 0, err
	}

	randomNumber := int(randomBigInt.Int64()) + minId
	return randomNumber, nil
}
