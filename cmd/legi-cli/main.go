package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/fabienogli/legigpt/api"
	"github.com/fabienogli/legigpt/httputils"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "legi-cli",
	Short: "Legi allows to search legi API",
	// Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := godotenv.Load(".env")
		if err != nil {
			return fmt.Errorf("when loading.env: %w", err)
		}

		clientID := os.Getenv("AIFE_CLIENT_ID")
		if clientID == "" {
			return fmt.Errorf("key AIFE_CLIENT_ID empty")
		}
		clientSecret := os.Getenv("AIFE_CLIENT_SECRET")
		if clientSecret == "" {
			return fmt.Errorf("key AIFE_CLIENT_SECRET empty")
		}
		ctx := cmd.Context()

		httpClient := httputils.NewResponseLsogger(
			httputils.NewClient(http.DefaultClient),
		)

		OauthCfg := api.OauthConfig{
			URL:          "https://sandbox-oauth.piste.gouv.fr/api/oauth/token",
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}

		tokenFile := path.Join(os.TempDir(), "token.json")
		log.Println("saving token into %s", tokenFile)
		fileStore := api.NewFileStore(tokenFile)

		authentifiedClient := api.NewOauthClient(OauthCfg, httpClient, fileStore)

		// tokenResponse, err := authentifiedClient.RetrievToken(ctx)

		// if err != nil {
		// 	log.Println(err)
		// 	return err
		// }
		// log.Println(tokenResponse)

		// to avoid searching for a token

		authClient := api.AuthentifiedClient{
			Client: authentifiedClient,
			URL:    "https://sandbox-api.piste.gouv.fr/dila/legifrance/lf-engine-app",
		}

		// results, err := authClient.Search(ctx, api.Search{
		// 	Recherche: api.Recherche{
		// 		Filtres: []api.Filtre{
		// 			// 	{
		// 			// 		Dates: Dates{
		// 			// 			Start: "2015-01-01",
		// 			// 			End:   "2023-01-31",
		// 			// 		},
		// 			// 		Facette: "DATE_SIGNATURE",
		// 			// 	},
		// 		},
		// 		Sort:                  "SIGNATURE_DATE_DESC",
		// 		FromAdvancedRecherche: false,
		// 		SecondSort:            "ID",
		// 		Champs: []api.Champ{
		// 			{
		// 				Criteres: []api.Critere{
		// 					{
		// 						Proximite: 2,
		// 						Valeur:    "dispositions",
		// 						Criteres: []api.Critere{
		// 							{
		// 								Valeur:        "maladie",
		// 								Operateur:     "ET",
		// 								TypeRecherche: "UN_DES_MOTS",
		// 							},
		// 							{
		// 								// Proximite:     3,
		// 								Valeur:        "congé",
		// 								Operateur:     "ET",
		// 								TypeRecherche: "UN_DES_MOTS",
		// 							},
		// 						},
		// 						Operateur:     "ET",
		// 						TypeRecherche: "UN_DES_MOTS",
		// 					},
		// 				},
		// 				Operateur: api.OperatorAND,
		// 				TypeChamp: api.FieldAll,
		// 			},
		// 		},
		// 		PageSize:       10,
		// 		Operateur:      api.OperatorAND,
		// 		TypePagination: api.PaginationDefault,
		// 		PageNumber:     1,
		// 	},
		// 	Fond: api.FondACCO,
		// })

		results, err := authClient.Consult(ctx, api.ConsultRequest{
			ID: "ACCOTEXT000037731479",
		})
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("results: %v", results.Acco.Attachment.Content)
		log.Println("results: %v", results.Acco.AttachementUrl)
		// err = authClient.Ping(ctx)
		// url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
		// oauth2.NewClient(context.Background(), oauth2.TokenSource{})
		return nil
	},
}

func init() {
	// rootCmd.PersistentFlags().Int64P("timeout", "t", 60, "duration in seconds that the command will be canceled")
	// rootCmd.PersistentFlags().Int64P("chunk-size", "s", 1, "the file will be split by this chunk size (MB)")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
