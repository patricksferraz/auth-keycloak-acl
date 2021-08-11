package external

import (
	"github.com/Nerzal/gocloak/v8"
)

type Keycloak struct {
	Realm        string
	ClientID     string
	ClientSecret string
	Audience     string
	Client       gocloak.GoCloak
}

func NewKeycloak(basePath, realm, clientID, clientSecret, audience string) *Keycloak {
	k := &Keycloak{
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
	}
	k.Client = gocloak.NewClient(basePath)

	return k
}
