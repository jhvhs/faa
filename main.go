package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"fmt"
	"strings"

	"github.com/concourse/faa/postfacto"
	"github.com/concourse/faa/slackcommand"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	retroPassword, _ := os.LookupEnv("POSTFACTO_RETRO_PASSWORD")

	retroApiURL, ok := os.LookupEnv("POSTFACTO_API_URL")
	if !ok {
		retroApiURL = "https://retros-iad-api.cfapps.io"
	}

	retroAppURL, ok := os.LookupEnv("POSTFACTO_APP_URL")
	if !ok {
		retroAppURL = "https://retros.cfapps.io"
	}

	vToken, ok := os.LookupEnv("SLACK_VERIFICATION_TOKEN")
	if !ok {
		panic(errors.New("must provide SLACK_VERIFICATION_TOKEN"))
	}

	retroID, ok := os.LookupEnv("POSTFACTO_RETRO_ID")
	if !ok {
		panic(errors.New("must provide POSTFACTO_RETRO_ID"))
	}

	c := &postfacto.RetroClient{
		ApiHost:  retroApiURL,
		AppHost: retroAppURL,
		ID:       retroID,
		Password: retroPassword,
	}

	server := slackcommand.Server{
		VerificationToken: vToken,
		Delegate: &PostfactoSlackDelegate{
			RetroClient: c,
		},
	}

	http.Handle("/", server)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type PostfactoSlackDelegate struct {
	RetroClient *postfacto.RetroClient
}

type Command string

const (
	CommandHappy Command = "happy"
	CommandMeh   Command = "meh"
	CommandSad   Command = "sad"
)

func (d *PostfactoSlackDelegate) Handle(r slackcommand.Command) (string, error) {
	parts := strings.SplitN(r.Text, " ", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("must be in the form of '%s [happy/meh/sad] [message]'", r.Command)
	}

	column := parts[0]
	description := parts[1]

	var (
		category postfacto.Category
	)

	switch Command(column) {
	case CommandHappy:
		category = postfacto.CategoryHappy
	case CommandMeh:
		category = postfacto.CategoryMeh
	case CommandSad:
		category = postfacto.CategorySad
	default:
		return "", errors.New("unknown command: must provide one of 'happy', 'meh' or 'sad'")
	}

	retroItem := postfacto.RetroItem{
		Category:    category,
		Description: fmt.Sprintf("%s [%s]", description, r.UserName),
	}

	err := d.RetroClient.Add(retroItem)
	if err != nil {
		return "", err
	}

	return "retro item added", nil
}
