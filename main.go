package main

import (
	"github.com/neosouler7/bitGoin/blockchain"
	"github.com/neosouler7/bitGoin/cli"
	"github.com/neosouler7/bitGoin/db"
)

func main() {
	defer db.Close()

	blockchain.Blockchain()
	cli.Start()

	// w := wallet.Wallet()
	// fmt.Println(w)
}
