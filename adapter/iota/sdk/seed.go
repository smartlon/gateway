package sdk

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//A Simple program using rand() to GEN an IOTA SEED for the IOTA GUI Wallet

const tokenLibrary = "ABCDEFGHIJKLMNOPQRSTUVXYWZ9"

func GenSeed() {
	keySize, _ := strconv.Atoi(os.Args[1])
	if keySize < 31 || keySize > 81 {
		fmt.Println("keySize must be grater than 30, lower than 82")
		os.Exit(1)
	}
	rand.Seed(time.Now().Unix())
	librarySize := len(tokenLibrary)
	for i:=0; i< keySize; i++ {
		p := rand.Intn(librarySize)
		fmt.Print(string(tokenLibrary[p]))
	}
	fmt.Println()
}