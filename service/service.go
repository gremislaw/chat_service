package service

import (
	"fmt"
	//"time"
	"os"
	"database/sql"
	"chat_service/api"
	"chat_service/core"
	"math/rand"
)

type ChatServiceServer struct {
	api.UnimplementedChatServiceServer
	db *sql.DB
}

type config struct {
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_DB       string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	APP_IP            string
	APP_PORT          string
}

func NewService(db *sql.DB) api.ChatServiceServer {
	service := ChatServiceServer{}

	/*for i := 1; !isDBAvailable(db); i++ {
		fmt.Printf("Db is unavailable(%ds)\n", i)
		time.Sleep(5 * time.Second)
	}*/
	fmt.Println("Db is available")
	service.db = db
	return &service
}

func isDBAvailable(db *sql.DB) bool {
	var res bool = true

	_, err := db.Query("select 1")

	if err != nil {
		res = false
	}
	return res
}

func NewConfig() config {
	cfg := config{
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		APP_IP:            os.Getenv("APP_IP"),
		APP_PORT:          os.Getenv("APP_PORT"),
	}
	return cfg
}

func (cfg config) IsValid() bool {
	var res bool = true
	if cfg.POSTGRES_USER == "" || cfg.POSTGRES_PASSWORD == "" ||
		cfg.POSTGRES_DB == "" || cfg.POSTGRES_HOST == "" ||
		cfg.POSTGRES_PORT == "" || cfg.APP_IP == "" ||
		cfg.APP_PORT == "" {
		res = false
	}
	return res
}

func (s *ChatServiceServer) HandleCommunication(stream api.ChatService_HandleCommunicationServer) error {

	clientUniqueCode := int64(rand.Intn(1e6))
	errch := make(chan error)

	go core.ReceiveFromStream(stream, clientUniqueCode, errch)
	go core.SendToStream(stream, clientUniqueCode, errch)

	return <-errch

}