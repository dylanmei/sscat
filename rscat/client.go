package rscat

import (
	"fmt"

	"github.com/rightscale/rsc/rsapi"
	"github.com/rightscale/rsc/ss/ssd"
	"github.com/rightscale/rsc/ss/ssm"
)

type Client struct {
	account  int
	designer *ssd.Api
	manager  *ssm.Api
}

func NewClient(apiHost, ssHost string, account int, token string) (*Client, error) {
	apiAuth := rsapi.NewOAuthAuthenticator(token, account)
	apiAuth.SetHost(apiHost)

	if err := apiAuth.CanAuthenticate(apiHost); err != nil {
		return nil, fmt.Errorf("invalid credentials: %s\n", err)
	}

	ssAuth := rsapi.NewSSAuthenticator(apiAuth, account)
	return &Client{
		account:  account,
		designer: ssd.New(ssHost, ssAuth),
		manager:  ssm.New(ssHost, ssAuth),
	}, nil
}
