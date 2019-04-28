package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"../config"
	"./helpers"
)

func main() {

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("Invalid arguments...")
	}

	if args[0] == "gen-token" {
		tokenGenerate(args[1:])
	} else {
		log.Fatal("Invalid action")
	}
}

func tokenGenerate(args []string) {
	// load application configurations
	tokenConfig := config.MustLoadConfig()

	aud := args[0]
	storeID, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Impossible to read store id")
	}
	storeName := args[2]
	expDays, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatal("Impossible to read expiration days")
	}

	helpers.SetupTokenHelper(tokenConfig.SessionToken.Secret,
		tokenConfig.SessionToken.Issuer,
		expDays,
		tokenConfig.SessionToken.Duration,
		tokenConfig.SessionToken.Duration,
	)

	token := helpers.Token.CreateStoreToken(aud, storeName, storeID)
	tokenString, err := helpers.Token.EncodeToken(token)
	if err != nil {
		log.Fatalf("Error found when creating store token. %s", err)
	}
	fmt.Printf("JWT Token:\n%s\n", tokenString)
}
