package main

import (
	"evm-downloader/http_client"
	"log"
)

func start() {
	http_client.Initialize()
	http_client.DownloadVscmrFiles()
}

func stop() {
	log.Printf("[+] finalizando programa")
}

func main() {
	// Inicio do programa
	start()

	// finalização
	defer stop()
}
