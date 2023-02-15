package main

import (
	cerberus "cerberus/src"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("<== Cerberus ==>")

	initCache()
	initDB()
	serve()

}

func initCache() {
	cacheHost := os.Getenv("CERBERUS_CACHE_HOST")
	cachePort := os.Getenv("CERBERUS_CACHE_PORT")

	fmt.Println("The cache would be expected at " + cacheHost + ":" + cachePort)
}

func initDB() {
	dbHost := os.Getenv("CERBERUS_DB_HOST")
	dbPort := os.Getenv("CERBERUS_DB_PORT")

	fmt.Println("The database would be expected at " + dbHost + ":" + dbPort)
}

func serve() {
	port := os.Getenv("CERBERUS_PORT")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Failed to listen on port " + port)
		panic(err)
	} else {
		fmt.Println("Listening on port " + port)
	}

	s := grpc.NewServer()
	reflection.Register(s) // added for services discovery
	cerberus.RegisterAuthenticationServer(s, &cerberus.Server{})
	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}
