package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/joaovds/learn-grpc-go/internal/db"
	"github.com/joaovds/learn-grpc-go/internal/pb"
	"github.com/joaovds/learn-grpc-go/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "50051"

func main() {
	database, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	authorDB := db.NewAuthor(database)
	authorService := service.NewAuthorService(*authorDB)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthorServiceServer(grpcServer, authorService)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panic("failed to listen: %w", err.Error())
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Panic("failed to serve: %w", err.Error())
	}
}
