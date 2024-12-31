package usecases

import (
	"api_crypto1.0/internal/db/repository"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type Usecases struct {
	Repository *repository.Repository
	URL        string
	Client     *ethclient.Client
}

func NewUsecases(rep *repository.Repository, url string) *Usecases {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	return &Usecases{
		Repository: rep,
		URL:        url,
		Client:     client,
	}
}
