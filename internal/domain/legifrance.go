package domain

type SearchQuery struct {
	Title string
	//PageSize Max=100
	//TODO: limit and offset
	PageSize   int
	PageNumber int
}

// TODO: nothing is named correctly
type AccordsWrapped struct {
	Accords []Accord
	Total   int
}

type Accord struct {
	ID    string
	CID   string
	Title string
}

type Content struct {
	Texte string
}
