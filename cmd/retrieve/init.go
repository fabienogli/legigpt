package retrieve

import (
	"log"
	"net/http"

	"github.com/fabienogli/legigpt"
	"github.com/fabienogli/legigpt/httputils"
	"github.com/fabienogli/legigpt/pkg/legifranceapi"
	"github.com/fabienogli/legigpt/pkg/store"
)

func initDealLooker(cfg legigpt.DealLookerConfiguration) *legifranceapi.AuthentifiedClient {
	httpClient :=
		// httputils.NewResponseLsogger(
		httputils.NewClient(http.DefaultClient)
	// )

	log.Println("saving token into %s", cfg.TokenFilename)

	fileStore := store.NewFileStore(cfg.TokenFilename)

	authentifiedClient := legifranceapi.NewOauthClient(cfg.LegiFranceConfiguration, httpClient, fileStore)

	return &legifranceapi.AuthentifiedClient{
		Client: authentifiedClient,
		URL:    "https://sandbox-api.piste.gouv.fr/dila/legifrance/lf-engine-app",
	}
}
