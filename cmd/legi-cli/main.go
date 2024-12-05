package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fabienogli/legigpt"
	"github.com/fabienogli/legigpt/httputils"
	"github.com/fabienogli/legigpt/internal/domain"
	"github.com/fabienogli/legigpt/internal/usecase"
	"github.com/fabienogli/legigpt/pkg/legifranceapi"
	"github.com/fabienogli/legigpt/pkg/store"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/mistral"
	"github.com/tmc/langchaingo/llms/ollama"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "legi-cli",
	Short: "Legi allows to search legi API",
	// Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := legigpt.InitConfiguration()
		if err != nil {
			return fmt.Errorf("when init Config: %w", err)
		}
		ctx := cmd.Context()

		authClient := initDealLooker(cfg.DealLookerConfiguration)

		dealLooker := usecase.NewLegifrance(authClient)

		llm, err := initGPT(cfg.GPTConfiguration)
		if err != nil {
			return fmt.Errorf("when llm new: %w", err)
		}

		gpt := usecase.NewGPT(llm, 40)

		top := usecase.NewDealLooker(dealLooker, gpt)

		err = top.Search(ctx, domain.SearchQuery{
			Title:      "congé",
			PageSize:   1,
			PageNumber: 1,
		})
		if err != nil {
			return fmt.Errorf("when dealLooker.Search: %w", err)
		}
		return nil
	},
}

func initGPT(cfg legigpt.GPTConfiguration) (llms.Model, error) {
	if cfg.Local != nil {
		return ollama.New(ollama.WithModel(*cfg.Local))

	}
	if cfg.Mistral != nil {
		return mistral.New(mistral.WithAPIKey(cfg.Mistral.ApiKey))
	}
	return nil, fmt.Errorf("when initializing gpt")
}

func initDealLooker(cfg legigpt.DealLookerConfiguration) *legifranceapi.AuthentifiedClient {
	httpClient := httputils.NewResponseLsogger(
		httputils.NewClient(http.DefaultClient),
	)

	log.Println("saving token into %s", cfg.TokenFilename)

	fileStore := store.NewFileStore(cfg.TokenFilename)

	authentifiedClient := legifranceapi.NewOauthClient(cfg.LegiFranceConfiguration, httpClient, fileStore)

	return &legifranceapi.AuthentifiedClient{
		Client: authentifiedClient,
		URL:    "https://sandbox-api.piste.gouv.fr/dila/legifrance/lf-engine-app",
	}
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
