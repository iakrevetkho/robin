package authgoogle

import (
	// External

	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	// Internal
	"github.com/iakrevetkho/robin/internal/config"
)

const (
	profileEndroint = "https://www.googleapis.com/oauth2/v2/userinfo"
)

type Provider struct {
	oAuthConf *oauth2.Config
	// URL which should be sent to user to auth via Google
	AuthURL string
}

func New(config config.Config) (provider *Provider, err error) {
	// Read JSON file with Google API secrets for your service
	data, err := ioutil.ReadFile(config.Auth.Google.SecretFileName)
	if err != nil {
		return
	}
	conf, err := google.ConfigFromJSON(data)
	if err != nil {
		return
	}
	// Add scopes for requesting data
	conf.Scopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"openid",
	}

	provider = &Provider{
		oAuthConf: conf,
		AuthURL:   conf.AuthCodeURL("state"),
	}
	return
}

// Method for processing redirect from auth provider after user is authenticated
func (p *Provider) ProcessAuthRedirect(authCode string) (err error) {
	// Handle the exchange code to initiate a transport.
	token, err := p.oAuthConf.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatal(err)
	}

	client := p.oAuthConf.Client(oauth2.NoContext, token)

	// Get user email from Google
	response, err := client.Get("https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%+v", response)
	return
}