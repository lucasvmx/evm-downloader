package model

type InfoBasica struct {
	Dg  string   `json:"dg"`
	Hg  string   `json:"hg"`
	Abr []Estado `json:"abr"`
}
