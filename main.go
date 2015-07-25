package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("rscat", "Author RightScale CAT files with less clicking")
	apiHost  = app.Flag("api.host", "RightScale API host name").Required().PlaceHolder("API_HOST").OverrideDefaultFromEnvar("RS_API_HOST").String()
	apiToken = app.Flag("api.token", "RightScale OAuth refresh token").Required().PlaceHolder("API_TOKEN").OverrideDefaultFromEnvar("RS_API_TOKEN").String()
	ssHost   = app.Flag("ss.host", "RightScale SelfService host name").Required().PlaceHolder("SS_HOST").OverrideDefaultFromEnvar("RS_SS_HOST").String()
	account  = app.Flag("account", "RightScale account number").Required().PlaceHolder("ACCOUNT").OverrideDefaultFromEnvar("RS_ACCOUNT").Int()
)

func main() {
	setup_upload_template(app)
	setup_delete_template(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
