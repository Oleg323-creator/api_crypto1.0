package usecases

import (
	"api_crypto1.0/internal/db/repository"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"os"
)

func (u *Usecases) Withdraw(addr string, amount string) (string, error) {

	fromAddress, err := u.Repository.GetRootAddr()
	if err != nil {
		return "", err
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

	privateKeyToDecription, nonc, err := u.Repository.GetRootPrivateKey(fromAddress)

	privateKey, err := u.DecryptAESGCM(nonc, privateKeyToDecription, password)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt private key for address %s: %v", fromAddress, err)
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

	log.Printf("Money is comming to your address! Hash: %s\n", signedTx.Hash().Hex())

	if err = u.SaveTxDataToDbByHash(signedTx.Hash().Hex(), bigIntAmount.String()); err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return signedTx.Hash().Hex(), nil
}

func (u *Usecases) SaveTxDataToDbByHash(hash string, amount string) error {

	tx, _, err := u.Client.TransactionByHash(context.Background(), common.HexToHash(hash))
	if err != nil {
		return err
	}

	chainID, err := u.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Error getting chain ID: %v", err)
	}

	signer := types.LatestSignerForChainID(chainID)
	senderAddr, err := types.Sender(signer, tx)
	if err != nil {
		return err
	}

	data := repository.Params{
		Hash:     tx.Hash().Hex(),
		FromAddr: senderAddr.String(),
		ToAddr:   tx.To().Hex(),
		Value:    amount,
		Currency: "ETH",
		TxType:   "Withdraw",
	}
	err = u.Repository.SaveTxDataToDB(data)
	if err != nil {
		log.Fatalf("Error saiving data in db: %v", err)
	}

	log.Println("Tx data saved to DB")

	return nil
}
