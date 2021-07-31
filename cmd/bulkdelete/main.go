package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

var (
	slackAPIToken string
	channelID     string
)

func getEnv() error {
	s, ok := os.LookupEnv("SLACK_API_TOKEN")
	if !ok {
		return errors.New("env SLACK_API_TOKEN is not found")
	}

	slackAPIToken = s

	c, ok := os.LookupEnv("CHANNEL_ID")
	if !ok {
		return errors.New("env CHANNEL_ID is not found")
	}

	channelID = c

	return nil
}

func main() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}

	api := slack.New(slackAPIToken)

	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     2,
	}

	res, err := api.GetConversationHistory(params)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range res.Messages {
		timestamp := msg.Timestamp

		a, b, err := api.DeleteMessage(channelID, timestamp)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("--------------- %+v\n%+v\n", a, b)
	}
}
