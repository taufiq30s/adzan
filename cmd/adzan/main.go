package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/taufiq30s/adzan/internal/routes"
)

func main() {
	defer catch()

	router := routes.NewRoute()

	port, err := getPort()
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}

}

func getPort() (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		return "", err
	}
	return os.Getenv("PORT"), nil
}

func catch() {
	if r := recover(); r != nil {
		log.Printf("panic: %v", r)
	}
}
