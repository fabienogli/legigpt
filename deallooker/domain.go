package deallooker

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
