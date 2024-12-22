package legigpt

import (
	"fmt"
	"os"
	"path"

	"github.com/fabienogli/legigpt/pkg/legifranceapi"
	"github.com/joho/godotenv"
)

type Configuration struct {
	DealLookerConfiguration DealLookerConfiguration
	GPTConfiguration        GPTConfiguration
	FolderStore             string
}

type DealLookerConfiguration struct {
	LegiFranceConfiguration legifranceapi.OauthConfig
	TokenFilename           string
}

type GPTConfiguration struct {
	Local   *string
	Mistral *GPTApi
}

type GPTApi struct {
	ApiKey string
}

func InitConfiguration() (Configuration, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return Configuration{}, fmt.Errorf("when loading.env: %w", err)
	}

	clientID := os.Getenv("AIFE_CLIENT_ID")
	if clientID == "" {
		return Configuration{}, fmt.Errorf("key AIFE_CLIENT_ID empty")
	}
	clientSecret := os.Getenv("AIFE_CLIENT_SECRET")
	if clientSecret == "" {
		return Configuration{}, fmt.Errorf("key AIFE_CLIENT_SECRET empty")
	}
	OauthCfg := legifranceapi.OauthConfig{
		URL:          "https://sandbox-oauth.piste.gouv.fr/api/oauth/token",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	folderStore := os.TempDir()

	// using mistral AI
	// not working
	// mistralAPIKEY := os.Getenv("MISTRAL_API_KEY")
	// if clientSecret == "" {
	// 	return Configuration{}, fmt.Errorf("key MISTRAL_API_KEY empty")
	// }
	gptLocal := "llama2"
	return Configuration{
		FolderStore: folderStore,
		DealLookerConfiguration: DealLookerConfiguration{
			LegiFranceConfiguration: OauthCfg,
			TokenFilename:           path.Join(folderStore, "token.json"),
		},
		GPTConfiguration: GPTConfiguration{
			Local: &gptLocal,
			// Mistral: &GPTApi{
			// 	ApiKey: mistralAPIKEY,
			// },
		},
	}, nil
}
