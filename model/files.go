package model

type Files struct {
	Dg     string
	Hg     string
	Hashes []FilesHash
}

type FilesHash struct {
	Hash  string
	Dr    string
	Hr    string
	St    string
	Nmarq []string
}

type ResourceInfo struct {
	ContentLength int64
}
