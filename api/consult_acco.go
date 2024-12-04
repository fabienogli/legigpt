package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ConsultRequest struct {
	SearchedString string `json:"searchedString"`
	ID             string `json:"id"`
}

// Main structure to represent the "acco" payload
type AccoPayload struct {
	Acco          AccoDetails `json:"acco"`
	ExecutionTime int         `json:"executionTime"`
	Dereferenced  bool        `json:"dereferenced"`
}

// Structure to represent the "acco" details
type AccoDetails struct {
	DateEffet                int64    `json:"dateEffet"`
	Themes                   []Theme  `json:"themes"`
	ConformeVersionIntegrale bool     `json:"conformeVersionIntegrale"`
	DateMaj                  int64    `json:"dateMaj"`
	Nature                   string   `json:"nature"`
	Signataires              []string `json:"signataires"`
	DateDiffusion            int64    `json:"dateDiffusion"`
	DateTexte                int64    `json:"dateTexte"`
	// RelevantDate             string          `json:"relevantDate"`
	Attachment          Attachment      `json:"attachment"`
	TitreTexte          string          `json:"titreTexte"`
	RefInjection        string          `json:"refInjection"`
	Url                 string          `json:"url"`
	Secteur             string          `json:"secteur"`
	ID                  string          `json:"id"`
	CodeIdcc            string          `json:"codeIdcc"`
	RaisonSociale       string          `json:"raisonSociale"`
	IDTechInjection     string          `json:"idTechInjection"`
	Origine             string          `json:"origine"`
	Numero              string          `json:"numero"`
	DateFin             int64           `json:"dateFin"`
	Syndicats           []Syndicat      `json:"syndicats"`
	AttachementUrl      string          `json:"attachementUrl"`
	CodeApe             string          `json:"codeApe"`
	AdressesPostales    []PostalAddress `json:"adressesPostales"`
	FileSize            string          `json:"fileSize"`
	DateDepot           int64           `json:"dateDepot"`
	CodeUniteSignataire string          `json:"codeUniteSignataire"`
	Data                string          `json:"data"`
	Siret               string          `json:"siret"`
}

// Structure for the "themes" array
type Theme struct {
	Libelle string `json:"libelle"`
	Code    string `json:"code"`
	Groupe  string `json:"groupe"`
}

// Structure for the "attachment" object
type Attachment struct {
	Title         string `json:"title"`
	Name          string `json:"name"`
	Language      string `json:"language"`
	Author        string `json:"author"`
	Keywords      string `json:"keywords"`
	Date          int64  `json:"date"`
	Content       string `json:"content"`
	ContentLength int64  `json:"content_length"`
	ContentType   string `json:"content_type"`
}

// Structure for the "syndicats" array
type Syndicat struct {
	Libelle string `json:"libelle"`
	Code    string `json:"code"`
}

// Structure for the "adressesPostales" array
type PostalAddress struct {
	Ville      string `json:"ville"`
	CodePostal string `json:"codePostal"`
}

// Partial implementation of search
func (a *AuthentifiedClient) Consult(ctx context.Context, search ConsultRequest) (AccoPayload, error) {
	payload, err := json.Marshal(search)
	if err != nil {
		return AccoPayload{}, fmt.Errorf("when json.Marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.URL+"/consult/acco", bytes.NewBuffer(payload))
	if err != nil {
		return AccoPayload{}, fmt.Errorf("whn http.NewRequest: %w", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	resp, err := a.Client.Do(ctx, req)
	if err != nil {
		return AccoPayload{}, fmt.Errorf("when Do: %w", err)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return AccoPayload{}, fmt.Errorf("when readAll: %w", err)
	}
	var results AccoPayload
	err = json.Unmarshal(buf, &results)
	if err != nil {
		return AccoPayload{}, fmt.Errorf("when json.Unmarshal: %w", err)
	}
	return results, nil
}
