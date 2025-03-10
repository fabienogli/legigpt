package retrieve

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/fabienogli/legigpt"
	"github.com/fabienogli/legigpt/internal/domain"
	"github.com/fabienogli/legigpt/internal/repository"
	"github.com/fabienogli/legigpt/internal/usecase"
	"github.com/fabienogli/legigpt/pkg/store"
	"github.com/spf13/cobra"
)

const (
	outputFlag = "output"
	searchFlag = "search"
)

const (
	defaultOutput = "history.json"
)

type appArg struct {
	output     string
	searchTerm string
}

func InitCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retrieve",
		Short: "Will retrieve the deals using the LegiAPI",
		// Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := legigpt.InitConfiguration()
			if err != nil {
				return fmt.Errorf("when init Config: %w", err)
			}

			outputFile, err := cmd.PersistentFlags().GetString(outputFlag)
			if err != nil {
				return fmt.Errorf("when getting output flag %w", err)
			}
			if outputFile == "" {
				return fmt.Errorf("output flag empty")
			}
			searchFlag, err := cmd.PersistentFlags().GetString(searchFlag)
			if err != nil {
				return fmt.Errorf("when getting search file: %w", err)
			}
			if searchFlag == "" {
				return fmt.Errorf("search flag empty")
			}
			//Storing file
			slog.Info("storing result", "file", outputFile)
			ctx := cmd.Context()
			return run(ctx, appArg{
				output:     outputFile,
				searchTerm: searchFlag,
			}, cfg.DealLookerConfiguration)
		},
	}
	folderStore := os.TempDir()
	outHistoryFile := path.Join(folderStore, defaultOutput)

	cmd.PersistentFlags().StringP(outputFlag, "o", outHistoryFile, "output file")
	cmd.PersistentFlags().StringP(searchFlag, "s", "télétravail", "search query for deals")

	return cmd
}

func run(ctx context.Context, arg appArg, cfg legigpt.DealLookerConfiguration) error {
	authClient := initDealLooker(cfg)

	dealLooker := usecase.NewLegifrance(authClient)

	top := usecase.NewDealLooker(dealLooker)

	query := domain.SearchQuery{
		Title:      arg.searchTerm,
		LimitSize:  40,
		PageNumber: 0,
	}

	slog.Info("looking for query", "query", query)

	deals, err := top.Search(ctx, query)
	if err != nil {
		return fmt.Errorf("when dealLooker.Search: %w", err)
	}

	historyFile := store.NewFileStore(arg.output)
	dbSearch := repository.NewDealRepository(historyFile)
	err = dbSearch.Store(ctx, domain.SearchHistory{
		Query:    query,
		Response: deals,
	})
	if err != nil {
		return fmt.Errorf("error storing searches: %w", err)
	}
	return nil
}
