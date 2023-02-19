package main

import (
	"database/sql"
	"strconv"

	"github.com/tortuga-softworks/cerberus/internal/authentication"
	"github.com/tortuga-softworks/cerberus/internal/session"
	"github.com/tortuga-softworks/cerberus/proto"

	"github.com/tortuga-softworks/hestia/pkg/account"

	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("<== Cerberus ==>")

	accountStore := initAccountStore()
	sessionStore := initSessionStore()                         // sessions storage management
	authService := initAuthService(accountStore, sessionStore) // auth logic
	authServer := initAuthServer(authService)                  // requests routing

	listener := initListener()

	server := grpc.NewServer()
	reflection.Register(server) // added for services discovery
	proto.RegisterAuthenticationServer(server, authServer)
	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}

func initAccountStore() account.AccountStore {
	dbConnectionString := os.Getenv("CERBERUS_ACCOUNTS_DB")

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL.")

	store, err := account.NewSqlAccountStore(db)

	if err != nil {
		panic(err)
	}

	fmt.Println("Account store ready.")
	return store
}

func initSessionStore() session.SessionStore {
	cacheHost := os.Getenv("CERBERUS_SESSIONS_HOST")
	cachePort := os.Getenv("CERBERUS_SESSIONS_PORT")

	var sessionDuration uint64

	sessionDurationString := os.Getenv("CERBERUS_SESSION_DURATION")
	if sessionDurationString == "" {
		sessionDuration = 43200
		fmt.Println("No session duration configration found. Using default: 43200.")
	} else {
		parsedSessionDuration, err := strconv.ParseUint(sessionDurationString, 10, 32)
		if err != nil {
			panic(err)
		}

		sessionDuration = parsedSessionDuration
	}

	sessionStore, err := session.NewRedisSessionStore(cacheHost+":"+cachePort, "", sessionDuration)

	if err != nil {
		panic(err)
	}

	fmt.Println("Session store ready.")
	return sessionStore
}

func initAuthService(accountStore account.AccountStore, sessionStore session.SessionStore) *authentication.AuthenticationService {
	authService, err := authentication.NewAuthenticationService(accountStore, sessionStore)

	if err != nil {
		panic(err)
	}

	fmt.Println("Authentication service ready.")
	return authService
}

func initAuthServer(authenticationService *authentication.AuthenticationService) *authentication.AuthenticationServer {
	authServer, err := authentication.NewAuthenticationServer(authenticationService)

	if err != nil {
		panic(err)
	}

	fmt.Println("Authentication server ready.")
	return authServer
}

func initListener() net.Listener {
	port := os.Getenv("CERBERUS_PORT")

	if port == "" {
		port = "9000"
		fmt.Println("No port configration found. Using default: 9000.")
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Listening on port " + port + ".")
	}

	return listener
}
