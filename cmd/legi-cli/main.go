package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

		authentifiedClient := api.NewOauthClient(OauthCfg, httpClient)

		authentifiedClient.Token = "ZXrDFTrDxQe86E8FhuWN5sFIIPvoifd61KfazDYQ1onxjC4c3Oaplb"

		authClient := api.AuthentifiedClient{
			Client: authentifiedClient,
			URL:    "https://sandbox-api.piste.gouv.fr/dila/legifrance/lf-engine-app",
		}
		// err = authClient.Search(ctx, "32h")
		err = authClient.Ping(ctx)
		log.Println(err)
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
