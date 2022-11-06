package main

import "evm-downloader/http_client"

func main() {
	http_client.Initialize()
	http_client.DownloadVscmrFiles()
}
