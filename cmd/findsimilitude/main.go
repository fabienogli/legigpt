package findsimilitude

import (
	"context"
	"fmt"
	"log/slog"
	"path"

	"github.com/fabienogli/legigpt"
	"github.com/fabienogli/legigpt/internal/repository"
	"github.com/fabienogli/legigpt/internal/usecase"
	"github.com/fabienogli/legigpt/pkg/store"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
)

const (
	readFromFile = "input"
)

type appArg struct {
	inputFile   string
	summarySize int
	search      string
}

func InitCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find-similitude",
		Short: "Will find similitude looking through deals powered by LLM",
		// Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := legigpt.InitConfiguration()
			if err != nil {
				return fmt.Errorf("when init Config: %w", err)
			}
			ctx := cmd.Context()

			llm, err := initGPT(cfg.GPTConfiguration)
			if err != nil {
				return fmt.Errorf("when llm new: %w", err)
			}

			outHistoryFile := path.Join(cfg.FolderStore, "history.json")
			return run(ctx, appArg{
				inputFile:   outHistoryFile,
				summarySize: 40,
				search:      "Je veux avoir autant de télétravail que je veux",
			}, llm)
		},
	}
	cmd.PersistentFlags().StringP(readFromFile, "i", "", "input file avoid calling web")
	cmd.MarkPersistentFlagRequired(readFromFile)
	return cmd
}

func run(ctx context.Context, arg appArg, llm llms.Model) error {
	gpt := usecase.NewGPT(llm, arg.summarySize)

	top := usecase.NewDealLLM(gpt)

	//Storing file
	slog.Info("storing result", "file", arg.inputFile)
	historyFile := store.NewFileStore(arg.inputFile)
	dbSearch := repository.NewDealRepository(historyFile)
	searchHistory, err := dbSearch.Get(ctx)

	if err != nil {
		return fmt.Errorf("error storing searches: %w", err)
	}

	slog.Info("Finding similarities", "based_search", searchHistory.Query)
	bestDeal, err := top.Rag(ctx, searchHistory.Response.Accords, arg.search)
	if err != nil {
		return fmt.Errorf("when rag: %w", err)
	}
	slog.Info("found best Deal", "best_deal", bestDeal)

	return nil
}
