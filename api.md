# Download de arquivos das urnas - TSE

## Presets

**Código 406** - eleições do primeiro turno 2022
**Código 407** - eleições do segundo turno de 2022
**String "ele2022"** - eleição de 2022


## Informações sobre municípios

`https://resultados.tse.jus.br/oficial/ele2022/544/config/mun-e000544-cm.json`

## 1 - Obter informações sobre as eleições

`https://resultados.tse.jus.br/oficial/comum/config/ele-c.json`

Processar os campos `ele` e `cd` - nome da eleição e código

Formato:
```json
{
	"dg": "01/11/2022",
	"hg": "17:11:27",
	"f": "O",
	"c": "ele2022",
	"pl": [
	    {
         "cd":"406",
         "cdpr":"395",
         "dt":"02/10/2022",
         "dtlim":"02/09/2024",
         "e": []
	    }
	]
```

Os outros campos podem ser ignorados

## 2 - Obter lista de seções e zonas eleitorais

`GET https://resultados.tse.jus.br/oficial/{ELE}/arquivo-urna/{CD}/config/df/df-p000{CD}-cs.json`

**CD - código obtido**
**ELE - nome da eleição obtido**

Exemplo:
```
https://resultados.tse.jus.br/oficial/ele2022/arquivo-urna/407/config/df/df-p000407-cs.json
```
## 3 - Nome do arquivo de boletim de urna

"o{CODIGO ELEICAO}-{CODIGO MUNICIPIO}{ZONA}{SECAO}.bu"

Exemplo:
`o00407-0112000080009.bu`

## 4 - Construir endpoint para baixar os arquivos de boletim de urna

Formato:
```
https://resultados.tse.jus.br/oficial/{ELE}/arquivo-urna/{CD}/dados/{NOME MUNICIPIO}/{CODIGO MUNICIPIO}/{ZONA}/{SECAO}/p000{CODIGO MUNICIPIO}-{NOME MUNICIPIO}-m{CODIGO MUNICIPIO}-z{ZONA}-s{SECAO}-aux.json
```

Exemplo:
`GET https://resultados.tse.jus.br/oficial/ele2022/arquivo-urna/407/dados/ac/01120/0008/0009/p000407-ac-m01120-z0008-s0009-aux.json`

Resposta esperada:

```json
{
   "dg":"31/10/2022",
   "hg":"00:45:20",
   "f":"O",
   "st":"Totalizada",
   "ds":"",
   "hashes":[
      {
         "hash":"77624b3969784a5a4e5276684942614f674779796a424b645a464d6a32795550684f683472785a494479733d",
         "dr":"30/10/2022",
         "hr":"17:19:30",
         "st":"Totalizado",
         "ds":"",
         "nmarq":[
            "o00407-0112000080009.imgbu",
            "o00407-0112000080009.vscmr",
            "o00407-0112000080009.logjez",
            "o00407-0112000080009.bu",
            "o00407-0112000080009.rdv"
         ]
      }
   ]
}
```
