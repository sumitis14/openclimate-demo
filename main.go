package main

import (
	// "encoding/json"
	"log"

	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/server"
	//"github.com/Varunram/essentials/ipfs"
	//"github.com/YaleOpenLab/openclimate/notif"
)

func main() {
	// Interact with the blockchain and check token balance
	blockchain.CheckTokenBalance()
	database.FlushDB()
	database.CreateHomeDir()
	log.Println("flusehd and created new db")
	server.StartServer("8001", true)
}
