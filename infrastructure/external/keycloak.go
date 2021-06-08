package external

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Nerzal/gocloak/v8"
	"github.com/joho/godotenv"
)

type Keycloak struct {
	BasePath     string
	Realm        string
	ClientID     string
	ClientSecret string
	Audience     string
	Client       gocloak.GoCloak
}

// TODO: Adds cover test
func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	if os.Getenv("ENV") == "dev" {
		err := godotenv.Load(basepath + "/../../.env")
		if err != nil {
			log.Fatal("Error loading .env files")
		}
	}

}

func ConnectKeycloak() *Keycloak {
	k := &Keycloak{
		BasePath:     os.Getenv("KEYCLOAK_BASE_PATH"),
		Realm:        os.Getenv("KEYCLOAK_REALM"),
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		Audience:     os.Getenv("KEYCLOAK_AUDIENCE"),
	}
	k.Client = gocloak.NewClient(k.BasePath)

	return k
}
