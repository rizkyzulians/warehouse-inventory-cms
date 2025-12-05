package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "password123"

	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password:", password)
	fmt.Println("Bcrypt Hash:", string(hash))
	fmt.Println("\nUpdate your seed data with this hash!")
}
