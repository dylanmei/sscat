package rscat

import (
	"fmt"

	"github.com/rightscale/rsc/rsapi"
	"github.com/rightscale/rsc/ss/ssd"
)

type Client struct {
	account int
	ssd     *ssd.Api
}

func NewClient(apiHost, ssHost string, account int, token string) (*Client, error) {
	auth := rsapi.NewOAuthAuthenticator(token, account)
	auth.SetHost(apiHost)

	if err := auth.CanAuthenticate(apiHost); err != nil {
		return nil, fmt.Errorf("invalid credentials: %s\n", err)
	}

	api := ssd.New(ssHost, rsapi.NewSSAuthenticator(auth, account))
	return &Client{account, api}, nil
}
