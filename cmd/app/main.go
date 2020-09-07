package main

import (
	"context"
	"github.com/DBoyara/Netology-Go-14/pkg/card"
	"github.com/DBoyara/Netology-Go-14/pkg/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	log.Printf("Server run on http://%s:%s", host, port)

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	dsn := "postgres://postgres:example@192.168.99.100:5432/db"
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Println(err)
		return
	}
	defer pool.Close()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Release()

	cardSvc := card.NewService()
	mux := http.NewServeMux()
	application := server.NewServer(cardSvc, mux, ctx, conn)
	application.Init()

	s := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return s.ListenAndServe()
}