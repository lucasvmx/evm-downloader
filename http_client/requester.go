package http_client

import (
	"encoding/json"
	"evm-downloader/model"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	fullDataSize  int64 = 0
	lastStateInfo []model.Estado
	lastState     string
)

func getTurno(turno int) (nomeTurno string) {
	if turno == model.SegundoTurno {
		return "407"
	}

	return "406"
}

func readBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[!] can't read HTTP body")
		return nil
	}

	return data
}

func compileZonesSectionsURL(nomeEleicao, turno, estado string) (url string) {
	var zonesSectionsURL string = "https://resultados.tse.jus.br/oficial/{ELEICAO}/arquivo-urna/{TURNO}/config/{ESTADO}/{ESTADO}-p000{TURNO}-cs.json"

	url = strings.ReplaceAll(zonesSectionsURL, "{ELEICAO}", nomeEleicao)
	url = strings.ReplaceAll(url, "{TURNO}", turno)
	url = strings.ReplaceAll(url, "{ESTADO}", estado)

	return
}

func compileURLForFile(nomeEleicao, turno, estado, cod_mun, zona, secao, hash, filename string) (url string) {
	var singleVscmrFileURL string = "https://resultados.tse.jus.br/oficial/{ELEICAO}/arquivo-urna/{TURNO}/dados/{ESTADO}/{CODIGO_MUNICIPIO}/{ZONA}/{SECAO}/{HASH}/{FILENAME}"

	url = strings.ReplaceAll(singleVscmrFileURL, "{ELEICAO}", nomeEleicao)
	url = strings.ReplaceAll(url, "{TURNO}", turno)
	url = strings.ReplaceAll(url, "{ESTADO}", estado)
	url = strings.ReplaceAll(url, "{CODIGO_MUNICIPIO}", cod_mun)
	url = strings.ReplaceAll(url, "{ZONA}", zona)
	url = strings.ReplaceAll(url, "{SECAO}", secao)
	url = strings.ReplaceAll(url, "{HASH}", hash)
	url = strings.ReplaceAll(url, "{FILENAME}", filename)

	return
}

func compileURLForfiles(nomeEleicao, turno, estado, cod_mun, sigla_mun, secao, zona string) (url string) {
	var urlBaseFiles string = "https://resultados.tse.jus.br/oficial/{ELEICAO}/arquivo-urna/{TURNO}/dados/{SIGLA_MUNICIPIO}/{CODIGO_MUNICIPIO}/{ZONA}/{SECAO}/p000{TURNO}-{SIGLA_MUNICIPIO}-m{CODIGO_MUNICIPIO}-z{ZONA}-s{SECAO}-aux.json"

	url = strings.ReplaceAll(urlBaseFiles, "{ELEICAO}", nomeEleicao)
	url = strings.ReplaceAll(url, "{TURNO}", turno)
	url = strings.ReplaceAll(url, "{ESTADO}", estado)
	url = strings.ReplaceAll(url, "{SIGLA_MUNICIPIO}", sigla_mun)
	url = strings.ReplaceAll(url, "{CODIGO_MUNICIPIO}", cod_mun)
	url = strings.ReplaceAll(url, "{ZONA}", zona)
	url = strings.ReplaceAll(url, "{SECAO}", secao)

	return
}

func PrintStatesInfo(basicInfo *model.InfoBasica) {
	states := basicInfo.Abr
	for _, state := range states {
		for _, mun := range state.Mu {
			log.Printf("==> %v", mun.Nm)
		}
	}
}

func DownloadBasicInfo() (basicInfo *model.InfoBasica) {
	var statesInfoURL string = "https://resultados.tse.jus.br/oficial/ele2022/544/config/mun-e000544-cm.json"

	resp, fail := http.Get(statesInfoURL)
	if fail != nil {
		log.Fatalf("[!] falha ao enviar GET para %v: %v", statesInfoURL, fail)
	}

	body := readBody(resp)
	if body == nil {
		log.Fatalf("[!] erro ao ler body da request")
	}

	fail = json.Unmarshal(body, &basicInfo)
	if fail != nil {
		log.Fatalf("[!] falha ao decodificar informacoes: %v", fail)
	}

	return
}

// DownloadZonesSectionsInfo baixa as informações das zonas eleitorais de um
func DownloadZonesSectionsInfo(state model.Estado, mun model.Municipio) (zones []model.Zona) {
	var info *model.InfoBasica

	if state.Cd == lastState {
		for _, _state := range lastStateInfo {
			for _, _mun := range _state.Mu {
				if _mun.Cd == mun.Cd {
					zones = _mun.Zon
					lastState = _state.Cd
					return
				}
			}
		}
	} else {
		log.Printf("[*] baixando dados do Estado: %v", state.Cd)
	}

	s := compileZonesSectionsURL("ele2022", getTurno(model.SegundoTurno), strings.ToLower(state.Cd))

	resp, err := http.Get(s)
	if err != nil {
		log.Fatalf("[!] erro ao baixar dados das secoes: %v", err)
	}

	data := readBody(resp)
	err = json.Unmarshal(data, &info)
	if err != nil {
		log.Fatalf("[!] erro ao decodificar body da requisicao")
	}

	lastStateInfo = info.Abr

	for _, _state := range info.Abr {
		for _, _mun := range _state.Mu {
			if _mun.Cd == mun.Cd {
				zones = _mun.Zon
				lastState = _state.Cd
				return
			}
		}
	}

	return
}

func getResourceInfo(url string) (info *model.ResourceInfo) {
	info = &model.ResourceInfo{
		ContentLength: 0,
	}

	resp, fail := http.Head(url)
	if fail != nil {
		log.Fatalf("[!] falha ao obter dados do recurso remoto: %v", fail)
	}

	contentLength := resp.Header.Get("Content-Length")
	v, _ := strconv.ParseInt(contentLength, 10, 64)
	info.ContentLength = v
	return
}

func getVscmrFileURL(fileURL, estado, cod_mun, zona, secao string) (remoteFileURL, localFilename string) {
	var files *model.Files

	customClient := &http.Client{}

	resp, fail := customClient.Get(fileURL)
	if fail != nil {
		log.Fatalf("[!] falha ao enviar request: %v", fail)
	}

	data := readBody(resp)
	fail = json.Unmarshal(data, &files)
	if fail != nil {
		log.Fatalf("[!] falha ao decodificar dados: %v (%v)", fail, fileURL)
	}

	hash := files.Hashes[0].Hash

	if len(files.Hashes[0].Nmarq) != 5 {
		return
	}

	filename := ""

	for _, file := range files.Hashes[0].Nmarq {
		if strings.Contains(file, ".vscmr") {
			filename = file
			break
		}
	}

	remoteFileURL = compileURLForFile("ele2022", getTurno(model.SegundoTurno), strings.ToLower(estado), cod_mun, zona, secao, hash, filename)
	localFilename = filename
	return
}

func downloadVscmrFile(fileURL, localFilename string) {
	customClient := &http.Client{}
	resp, fail := customClient.Get(fileURL)
	if fail != nil {
		log.Fatalf("[!] erro ao enviar request: %v", fail)
	}

	fileData := readBody(resp)
	os.WriteFile(localFilename, fileData, os.ModePerm)

	// Obtém o tamanho do download a ser realizado
	resInfo := getResourceInfo(fileURL)
	fullDataSize += resInfo.ContentLength
}

func DownloadVscmrFiles() {

	// Realiza o mapeamento dos arquivos
	CreateURLMap()

	// Faz o Download para o local atual
	for key, url := range urlMap {
		localFile := strings.Split(key, "|")[1]
		downloadVscmrFile(url, localFile)

		log.Printf("[*] baixado: %.2f KB", float64(fullDataSize/1024.0))
	}
}
