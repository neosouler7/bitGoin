package main

import (
	"fmt"

	"github.com/neosouler7/bitGoin/rest"
)

func main() {
	fmt.Println("Welcome to bitGoin :)")

	// chain := blockchain.GetBlockchain()
	// chain.AddBlock("Second Block")

	// go explorer.Start(3000)

	rest.Start(4000)
}
