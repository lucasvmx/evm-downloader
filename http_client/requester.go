package http_client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"urna-downloader/model"
)

var (
	urlBaseFiles       = "https://resultados.tse.jus.br/oficial/{ELEICAO}/arquivo-urna/{TURNO}/dados/{SIGLA_MUNICIPIO}/{CODIGO_MUNICIPIO}/{ZONA}/{SECAO}/p000{TURNO}-{SIGLA_MUNICIPIO}-m{CODIGO_MUNICIPIO}-z{ZONA}-s{SECAO}-aux.json"
	zonesSectionsURL   = "https://resultados.tse.jus.br/oficial/{ELEICAO}/arquivo-urna/{TURNO}/config/{ESTADO}/{ESTADO}-p000{TURNO}-cs.json"
	singleVscmrFileURL = "https://resultados.tse.jus.br/oficial/{ELEICAO}/arquivo-urna/{TURNO}/dados/{ESTADO}/{CODIGO_MUNICIPIO}/{ZONA}/{SECAO}/{HASH}/{FILENAME}"
)

var (
	fullDataSize int64 = 0
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
	url = strings.ReplaceAll(zonesSectionsURL, "{ELEICAO}", nomeEleicao)
	url = strings.ReplaceAll(url, "{TURNO}", turno)
	url = strings.ReplaceAll(url, "{ESTADO}", estado)

	log.Printf("[*] compiled url: %v", url)
	return
}

func compileURLForFile(nomeEleicao, turno, estado, cod_mun, zona, secao, hash, filename string) (url string) {
	url = strings.ReplaceAll(singleVscmrFileURL, "{ELEICAO}", nomeEleicao)
	url = strings.ReplaceAll(url, "{TURNO}", turno)
	url = strings.ReplaceAll(url, "{ESTADO}", estado)
	url = strings.ReplaceAll(url, "{CODIGO_MUNICIPIO}", cod_mun)
	url = strings.ReplaceAll(url, "{ZONA}", zona)
	url = strings.ReplaceAll(url, "{SECAO}", secao)
	url = strings.ReplaceAll(url, "{HASH}", hash)
	url = strings.ReplaceAll(url, "{FILENAME}", filename)

	log.Printf("[*] compiled url (single file): %v", url)
	return
}

func compileURLForfiles(nomeEleicao, turno, estado, cod_mun, sigla_mun, secao, zona string) (url string) {

	url = strings.ReplaceAll(urlBaseFiles, "{ELEICAO}", nomeEleicao)
	url = strings.ReplaceAll(url, "{TURNO}", turno)
	url = strings.ReplaceAll(url, "{ESTADO}", estado)
	url = strings.ReplaceAll(url, "{SIGLA_MUNICIPIO}", sigla_mun)
	url = strings.ReplaceAll(url, "{CODIGO_MUNICIPIO}", cod_mun)
	url = strings.ReplaceAll(url, "{ZONA}", zona)
	url = strings.ReplaceAll(url, "{SECAO}", secao)

	//log.Printf("[*] compiled url (for files): %v", url)

	return
}

func PrintStatesInfo(basicInfo *model.InfoBasica) {
	states := basicInfo.Abr
	for _, state := range states {
		log.Printf("Estado: %v (%v)", state.Cd, state.Ds)
		for _, mun := range state.Mu {
			log.Printf("==> %v", mun.Nm)
		}
	}
}

func DownloadBasicInfo() (basicInfo *model.InfoBasica) {
	var statesInfoURL string = "https://resultados.tse.jus.br/oficial/ele2022/544/config/mun-e000544-cm.json"

	resp, fail := http.Get(statesInfoURL)
	if fail != nil {
		log.Fatalf("[!] failed to send GET for %v: %v", statesInfoURL, fail)
	}

	body := readBody(resp)
	if body == nil {
		log.Fatalf("[!] can't read body")
	}

	fail = json.Unmarshal(body, &basicInfo)
	if fail != nil {
		log.Fatalf("[!] couldn't decode states information: %v", fail)
	}

	return
}

// DownloadZonesSectionsInfo baixa as informações das zonas eleitorais de um
func DownloadZonesSectionsInfo(state model.Estado, mun model.Municipio) (zones []model.Zona) {
	var info *model.InfoBasica

	s := compileZonesSectionsURL("ele2022", getTurno(model.SegundoTurno), strings.ToLower(state.Cd))

	resp, err := http.Get(s)
	if err != nil {
		log.Fatalf("[!] failed to download zones sections information: %v", err)
	}

	data := readBody(resp)
	json.Unmarshal(data, &info)

	for _, q := range info.Abr {
		for _, s := range q.Mu {
			if s.Cd == mun.Cd {
				zones = s.Zon
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
		log.Fatalf("[!] failed to get resource info")
	}

	contentLength := resp.Header.Get("Content-Length")
	log.Printf("len: %v", contentLength)
	v, _ := strconv.ParseInt(contentLength, 10, 64)
	info.ContentLength = v
	return
}

func downloadVscmrFile(fileURL, estado, cod_mun, zona, secao string) {
	var files *model.Files

	resp, fail := http.Get(fileURL)
	if fail != nil {
		log.Fatalf("[!] failed to send request: %v", fail)
	}

	data := readBody(resp)
	fail = json.Unmarshal(data, &files)
	if fail != nil {
		log.Fatalf("[!] failed to decode data: %v", fail)
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

	remoteFileURL := compileURLForFile("ele2022", getTurno(model.SegundoTurno), strings.ToLower(estado), cod_mun, zona, secao, hash, filename)

	resInfo := getResourceInfo(remoteFileURL)
	log.Printf("[*] downloading %v bytes of data", resInfo.ContentLength)

	fullDataSize += resInfo.ContentLength
}

func DownloadVscmrFiles() {

	// Faz o download das informações básicas dos Estados
	basicInfo := DownloadBasicInfo()

	// Para cada estado, listar os municípios e para cada município listar as zonas eleitorais
	for _, state := range basicInfo.Abr {
		for _, mun := range state.Mu {
			log.Printf("[+] processando seções e zonas de %v - %v ...", mun.Nm, state.Ds)

			zonesList := DownloadZonesSectionsInfo(state, mun)

			for _, zone := range zonesList {
				for _, section := range zone.Sec {
					url := compileURLForfiles("ele2022", getTurno(model.SegundoTurno), strings.ToLower(state.Cd), mun.Cd, strings.ToLower(state.Cd), section.Ns, zone.Cd)
					downloadVscmrFile(url, strings.ToLower(state.Cd), mun.Cd, zone.Cd, section.Ns)
				}
			}
		}
	}

	log.Printf("Tamanho total dos dados: %v KBytes", fullDataSize/1024)
}
