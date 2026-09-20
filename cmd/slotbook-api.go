package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/Aman-Shitta/slotbook/internal/server"
	"github.com/joho/godotenv"
)

//go:embed logo.txt
var banner string
var EnvPath string = ".env"

func main() {
	fmt.Println(banner)
	// load env

	if err := godotenv.Load(EnvPath); err != nil {
		log.Fatal("Something went Wrong")
	}

	var port string = os.Getenv("SLOTBOOK_ADDR")

	server := server.SetupRouter(port)

	server.Run(fmt.Sprintf(":%s", port))

}
