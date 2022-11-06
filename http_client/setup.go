package http_client

import "sync"

var (
	urlMap        map[string]string
	qtdMunicipios = 6000
	threadsChan   chan int
	maxThreads    = 6
	globalMux     *sync.Mutex
)

func Initialize() {
	urlMap = make(map[string]string, qtdMunicipios)
	threadsChan = make(chan int, maxThreads)
	globalMux = &sync.Mutex{}
}
