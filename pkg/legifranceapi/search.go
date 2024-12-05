package legifranceapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/fabienogli/legigpt/httputils"
)

type Search struct {
	Recherche `json:"recherche"`
	Fond      Fond `json:"fond"`
}

type Recherche struct {
	Filtres               []Filtre `json:"filtres"`
	Sort                  string   `json:"sort"`
	FromAdvancedRecherche bool     `json:"fromAdvancedRecherche"`
	SecondSort            string   `json:"secondSort"`
	Champs                []Champ  `json:"champs"`
	//PageSize max: (max=100)
	PageSize       int        `json:"pageSize"`
	Operateur      Operator   `json:"operateur"`
	TypePagination Pagination `json:"typePagination"`
	PageNumber     int        `json:"pageNumber"`
}

type Filtre struct {
	Dates   Dates  `json:"dates"`
	Facette string `json:"facette"`
}

type Dates struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Champ struct {
	Criteres  []Critere `json:"criteres"`
	Operateur Operator  `json:"operateur"`
	TypeChamp FieldType `json:"typeChamp"`
}

type Critere struct {
	Proximite     any          `json:"proximite"`
	Valeur        string       `json:"valeur"`
	Criteres      []Critere    `json:"criteres,omitempty"`
	Operateur     Operator     `json:"operateur"`
	TypeRecherche SearchedType `json:"typeRecherche"`
}

// SearchedType represents the type of search criteria or fields
type SearchedType string

const (
	// Enum values for SearchedType
	TypeUnDesMots                         SearchedType = "UN_DES_MOTS"
	TypeExacte                            SearchedType = "EXACTE"
	TypeTousLesMotsDansUnChamp            SearchedType = "TOUS_LES_MOTS_DANS_UN_CHAMP"
	TypeAucunDesMots                      SearchedType = "AUCUN_DES_MOTS"
	TypeAucuneCorrespondanceCetExpression SearchedType = "AUCUNE_CORRESPONDANCE_A_CETTE_EXPRESSION"
)

// FieldType represents the type of search fields
type FieldType string

const (
	// Enum values for FieldType
	FieldAll             FieldType = "ALL"
	FieldTitle           FieldType = "TITLE"
	FieldTable           FieldType = "TABLE"
	FieldNor             FieldType = "NOR"
	FieldNum             FieldType = "NUM"
	FieldAdvancedTexteID FieldType = "ADVANCED_TEXTE_ID"
	FieldNumDelib        FieldType = "NUM_DELIB"
	FieldNumDec          FieldType = "NUM_DEC"
	FieldNumArticle      FieldType = "NUM_ARTICLE"
	FieldArticle         FieldType = "ARTICLE"
	FieldMinistere       FieldType = "MINISTERE"
	FieldVisa            FieldType = "VISA"
	FieldNotice          FieldType = "NOTICE"
	FieldVisaNotice      FieldType = "VISA_NOTICE"
	FieldTravauxPrep     FieldType = "TRAVAUX_PREP"
	FieldSignature       FieldType = "SIGNATURE"
	FieldNota            FieldType = "NOTA"
	FieldNumAffaire      FieldType = "NUM_AFFAIRE"
	FieldAbstrats        FieldType = "ABSTRATS"
	FieldResumes         FieldType = "RESUMES"
	FieldTexte           FieldType = "TEXTE"
	FieldECLI            FieldType = "ECLI"
	FieldNumLoiDef       FieldType = "NUM_LOI_DEF"
	FieldTypeDecision    FieldType = "TYPE_DECISION"
	FieldNumeroInterne   FieldType = "NUMERO_INTERNE"
	FieldRefPubli        FieldType = "REF_PUBLI"
	FieldResumeCirc      FieldType = "RESUME_CIRC"
	FieldTexteRef        FieldType = "TEXTE_REF"
	FieldTitreLoiDef     FieldType = "TITRE_LOI_DEF"
	FieldRaisonSociale   FieldType = "RAISON_SOCIALE"
	FieldMotsCles        FieldType = "MOTS_CLES"
	FieldIDCC            FieldType = "IDCC"
)

type Operator string

const (
	OperatorOr  Operator = "OU"
	OperatorAND Operator = "ET"
)

type Pagination string

const (
	PaginationDefault Pagination = "DEFAUT"
	PaginationArticle Pagination = "ARTICLE"
)

// Fond defines the string-based enum for the given values with the prefix "Fond"
type Fond string

// Enum values for Fond with the "Fond" prefix
const (
	FondJORF      Fond = "JORF"
	FondCNIL      Fond = "CNIL"
	FondCETAT     Fond = "CETAT"
	FondJURI      Fond = "JURI"
	FondJUFI      Fond = "JUFI"
	FondCONSTIT   Fond = "CONSTIT"
	FondKALI      Fond = "KALI"
	FondCODE_DATE Fond = "CODE_DATE"
	FondCODE_ETAT Fond = "CODE_ETAT"
	FondLODA_DATE Fond = "LODA_DATE"
	FondLODA_ETAT Fond = "LODA_ETAT"
	FondALL       Fond = "ALL"
	FondCIRC      Fond = "CIRC"
	FondACCO      Fond = "ACCO"
)

type Response struct {
	ExecutionTime     int      `json:"executionTime"`
	Results           []Result `json:"results"`
	TotalResultNumber int      `json:"totalResultNumber"`
}

type Result struct {
	Titles []Title `json:"titles"`
}

type Title struct {
	ID    string `json:"id"`
	CID   string `json:"cid"`
	Title string `json:"title"`
}

type AuthentifiedClient struct {
	Client httputils.Doer
	URL    string
}

func (r *Response) String() string {
	s := ""
	for _, result := range r.Results {
		for _, title := range result.Titles {
			s += fmt.Sprintf(`
ID: %s
CID: %s
Title: %s
`, title.ID, title.CID, title.Title)
		}
	}
	return fmt.Sprintf(`
Titles: %s
TotalResultNumber: %d
`, s, r.TotalResultNumber)
}

// Partial implementation of search
func (a *AuthentifiedClient) Search(ctx context.Context, search Search) (Response, error) {
	payload, err := json.Marshal(search)
	if err != nil {
		return Response{}, fmt.Errorf("when json.Marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.URL+"/search", bytes.NewBuffer(payload))
	if err != nil {
		return Response{}, fmt.Errorf("whn http.NewRequest: %w", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	resp, err := a.Client.Do(ctx, req)
	if err != nil {
		return Response{}, fmt.Errorf("when Do: %w", err)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("when readAll: %w", err)
	}
	var results Response
	err = json.Unmarshal(buf, &results)
	if err != nil {
		slog.Info("error unmarshall", "status", resp.Status, "body", string(buf))
		return Response{}, fmt.Errorf("when json.Unmarshal: %w", err)
	}
	return results, nil
}
