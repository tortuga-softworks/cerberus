package main

import (
	"cerberus/internal/authentication"
	"cerberus/internal/session"
	pb "cerberus/proto"
	"strconv"

	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("<== Cerberus ==>")

	sessionStore := initSessionStore()                                                    // sessions storage management
	authService := initAuthService(sessionStore)                                          // auth logic
	authServer := authentication.AuthenticationServer{AuthenticationService: authService} // requests routing

	listener := initListener()

	server := grpc.NewServer()
	reflection.Register(server) // added for services discovery
	pb.RegisterAuthenticationServer(server, &authServer)
	if err := server.Serve(listener); err != nil {
		panic(err)
	}

}

func initSessionStore() session.SessionStore {
	cacheHost := os.Getenv("CERBERUS_SESSIONS_HOST")
	cachePort := os.Getenv("CERBERUS_SESSIONS_PORT")

	var sessionDuration int

	sessionDurationString := os.Getenv("CERBERUS_SESSION_DURATION")
	if sessionDurationString == "" {
		sessionDuration = 43200
		fmt.Println("No session duration configration found. Using default: 43200.")
	}

	sessionDuration, err := strconv.Atoi(sessionDurationString)
	if err != nil {
		panic(err)
	}

	sessionStore, err := session.NewRedisSessionStore(cacheHost+":"+cachePort, "", sessionDuration)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Session store connection established at " + cacheHost + ":" + cachePort)
	}

	return sessionStore
}

func initAuthService(sessionStore session.SessionStore) *authentication.AuthenticationService {
	authService, err := authentication.NewAuthenticationService(sessionStore)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Authentication service ready.")
	}

	return authService
}

func initListener() net.Listener {
	port := os.Getenv("CERBERUS_PORT")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Listening on port " + port + ".")
	}

	return listener
}
