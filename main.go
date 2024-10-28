package main

import (
	"fmt"
	//"database/sql"
	"google.golang.org/grpc"
	"chat_service/service"
	"net"
	"chat_service/api"
)

func main() {
	cfg := service.NewConfig()
	if !cfg.IsValid() {
		panic("incorrect config")
	}

	/*connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD, cfg.POSTGRES_DB, cfg.POSTGRES_HOST, cfg.POSTGRES_PORT)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("DB connection error")
	}
	defer db.Close()*/

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.APP_IP, cfg.APP_PORT))
	if err != nil {
		panic("listen error")
	}
	api.RegisterChatServiceServer(grpcServer, service.NewService(nil))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}