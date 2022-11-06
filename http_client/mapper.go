package http_client

import (
	"evm-downloader/filesystem"
	"evm-downloader/model"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func addItemToMap(zoneCd, sectionNs, munCd, stateCd string) {
	url := compileURLForfiles("ele2022", getTurno(model.SegundoTurno), strings.ToLower(stateCd), munCd, strings.ToLower(stateCd), sectionNs, zoneCd)
	vscmrFileURL, localFilename := getVscmrFileURL(url, strings.ToLower(stateCd), munCd, zoneCd, sectionNs)
	keyName := fmt.Sprintf("%v%v%v%v|%v", zoneCd, sectionNs, munCd, stateCd, localFilename)

	urlChannel <- fmt.Sprintf("%v|%v", vscmrFileURL, localFilename)

	globalMux.Lock()
	urlMap[keyName] = vscmrFileURL
	globalMux.Unlock()

	<-threadsChan
}

func saveURLMapToDisk() {

	var flags int = os.O_WRONLY

	now := time.Now()
	filename := fmt.Sprintf("url_maps_%v%v%v.dat", now.Day(), now.Month(), now.Year())

	if filesystem.FileExists(filename) {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_CREATE
	}

	filePtr, err := os.OpenFile(filename, flags, os.ModePerm)
	if err != nil {
		log.Printf("[!] erro ao criar arquivo de mapas: %v", err)
		return
	}

	defer filePtr.Close()

	for _, url := range urlMap {
		filePtr.WriteString(fmt.Sprintf("%v\n", url))
	}
}

func CreateURLMap() {
	// Faz o download das informações básicas dos Estados
	basicInfo := DownloadBasicInfo()

	// Para cada estado, listar os municípios e para cada município listar as zonas eleitorais
	for _, state := range basicInfo.Abr {
		log.Printf("[*] mapeando cidades do %v", state.Cd)
		for _, mun := range state.Mu {
			zonesList := DownloadZonesSectionsInfo(state, mun)

			for _, zone := range zonesList {
				for i, section := range zone.Sec {
					threadsChan <- i
					go addItemToMap(zone.Cd, section.Ns, mun.Cd, state.Cd)
				}
			}

			time.Sleep(time.Second * 5)
			saveURLMapToDisk()
		}
	}

	progressMux.Lock()
	totalFilesToDownload = len(urlMap)
	progressMux.Unlock()
	log.Printf("[*] quantidade url's mapeadas: %v", len(urlMap))
}
