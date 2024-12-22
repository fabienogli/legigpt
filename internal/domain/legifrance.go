package domain

type SearchQuery struct {
	Title string `json:"title"`
	//LimitSize Max=100
	//TODO: limit and offset
	LimitSize  int `json:"limit_size"`
	PageNumber int `json:"PageNumber"`
}

// TODO: nothing is named correctly
type DealResult struct {
	Accords []Accord `json:"accords"`
	Total   int      `json:"total"`
}

type Accord struct {
	ID    string `json:"id"`
	CID   string
	Title string `json:"title"`
	Texte string `json:"texte"`
}

type SearchHistory struct {
	Query    SearchQuery `json:"query"`
	Response DealResult  `json:"response"`
}
