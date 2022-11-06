package http_client

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func tasksCompleted() bool {
	progressMux.Lock()
	defer progressMux.Unlock()

	if totalFilesToDownload != -1 && (downloadedFiles == totalFilesToDownload) {
		return true
	}

	return false
}

func ShowProgress() {

	for {
		if tasksCompleted() {
			break
		}

		log.Printf("[*] baixado: %.2f MB", float64(fullDataSize/1024.0/1024.0))
		time.Sleep(time.Second * 20)
	}

	log.Printf("[+] finalizando thread de progresso")
}

// PerformDownload realiza o download de arquivos passados por um canal na memória
func PerformDownload() {

	var folderName string = "output/vscmr"
	os.MkdirAll(folderName, os.ModePerm|os.ModeDir)

	go ShowProgress()

	for {
		content := <-urlChannel
		contents := strings.Split(content, "|")
		url := contents[0]
		localFilename := contents[1]

		// Faz o Download para o local atual
		downloadVscmrFile(url, fmt.Sprintf("%v/%v", folderName, localFilename))

		downloadedFiles++
		time.Sleep(time.Second)

		if tasksCompleted() {
			break
		}
	}

	log.Printf("[+] finalizando thread de download")
}
