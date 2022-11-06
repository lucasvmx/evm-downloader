package http_client

import "sync"

var (
	urlMap               map[string]string
	qtdMunicipios        = 6000
	threadsChan          chan int
	maxThreads           = 6
	maxPendingDownloads  = 20
	globalMux            *sync.Mutex
	progressMux          *sync.Mutex
	urlChannel           chan string
	totalFilesToDownload int = -1
	downloadedFiles      int = 0
)

func Initialize() {
	urlMap = make(map[string]string, qtdMunicipios)
	threadsChan = make(chan int, maxThreads)
	urlChannel = make(chan string, maxPendingDownloads)
	globalMux = &sync.Mutex{}
	progressMux = &sync.Mutex{}
}
