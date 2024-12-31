package main

import (
	"api_crypto1.0/internal/db"
	"api_crypto1.0/internal/db/repository"
	"api_crypto1.0/internal/handlers"
	"api_crypto1.0/internal/runners"
	"api_crypto1.0/internal/usecases"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := db.ConnectionConfig{
		Host:     "localhost",
		Port:     "5435",
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	dbConn := db.NewDB(cfg)

	err = db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	rep := repository.NewRepository(dbConn)

	url := fmt.Sprintf("wss://sepolia.infura.io/ws/v3/%s", os.Getenv("API_KEY"))
	if os.Getenv("API_KEY") == "" {
		log.Fatal("Error getting api key from .env")
	}

	usecase := usecases.NewUsecases(rep, url)

	_, err = usecase.CheckRootAddrInDB()
	if err != nil {
		log.Fatalf("Error checking root address(MAIN) :%v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	runner := runners.NewRunner(usecase, ctx, wg)

	err = runner.RestartListeningAllAddr()
	if err != nil {
		log.Fatalf("Error listening adresses: %v", err)
	}

	handler := handlers.NewHandler(runner)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Logger(), handler.PostRoutineMiddleware())
	router.Use(gin.Recovery())

	router.POST("/create/address", handler.GetNewAddr)

	go func() {
		err = router.Run(":8080")
		if err != nil {
			log.Fatal("Failed to start Gin server:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("after stop")

	cancel()
	log.Println("after cancel")
	wg.Wait()
	log.Println("after Wait")
	fmt.Println("is ok")

}
