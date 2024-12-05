package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/fabienogli/legigpt/api"
	"github.com/fabienogli/legigpt/deallooker"
	"github.com/fabienogli/legigpt/httputils"
	"github.com/fabienogli/legigpt/llmx"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms/ollama"
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

		authClient := &api.AuthentifiedClient{
			Client: authentifiedClient,
			URL:    "https://sandbox-api.piste.gouv.fr/dila/legifrance/lf-engine-app",
		}

		dealLooker := deallooker.NewDealLooker(authClient)

		//local
		llm, err := ollama.New(ollama.WithModel("smollm"))
		if err != nil {
			return fmt.Errorf("")
		}

		//using mistral AI
		// not working
		// mistralAPIKEY := os.Getenv("MISTRAL_API_KEY")
		// if clientSecret == "" {
		// 	return fmt.Errorf("key MISTRAL_API_KEY empty")
		// }
		// llm, err := mistral.New(mistral.WithAPIKey(mistralAPIKEY))
		if err != nil {
			return fmt.Errorf("when llm new: %w", err)
		}

		gpt := llmx.NewGPT(llm)

		// summary, err := gpt.Summarize(ctx, "petit texte à résumé")
		// if err != nil {
		// 	return err
		// }
		// log.Println(summary)
		// return nil

		accords, err := dealLooker.Search(ctx, deallooker.SearchQuery{
			Title:      "congé",
			PageSize:   1,
			PageNumber: 1,
		})
		if err != nil {
			return fmt.Errorf("when dealLooker.Search: %w", err)
		}
		log.Println("total: ", accords.Total)

		var contents []deallooker.Content
		for _, accord := range accords.Accords {
			content, err := dealLooker.GetContent(ctx, accord.ID)
			if err != nil {
				return fmt.Errorf("when dealLooker.GetContent: %w", err)
			}
			summary, err := gpt.Summarize(ctx, content.Texte)
			if err != nil {
				return fmt.Errorf("when ollamaTest: %w", err)
			}
			log.Println("Summary: ", summary)
			contents = append(contents, content)
		}
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
