package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

func bulkDelete(api *slack.Client, oldestTimeStamp, cursor string) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     2,
		Oldest:    oldestTimeStamp,
		Cursor:    cursor,
	}

	res, err := api.GetConversationHistory(params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("--------------- %+v\n", res.Messages[0])

	for _, msg := range res.Messages {
		a, b, err := api.DeleteMessage(channelID, msg.Timestamp)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("--------------- %+v\n%+v\n", a, b)
	}

	if res.HasMore {
		bulkDelete(api, oldestTimeStamp, res.ResponseMetaData.NextCursor)
	}
}

func main() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}

	const shortForm = "2006/01/02"
	var (
		s string
		t time.Time
	)

	fmt.Println("This program will delete SLACK messages older than the entered date.")
	fmt.Printf("Enter date in the format like %s. :  ", shortForm)

	n, err := fmt.Scanln(&s)

	switch {
	case n == 0 && err.Error() == "unexpected newline":
		t = time.Now().AddDate(0, -1, 0)
	case err != nil:
		log.Fatal(err)
	default:
		t, err = time.Parse(shortForm, s)
		if err != nil {
			log.Fatal(err)
		}
	}

	oldestTimeStamp := strconv.FormatInt(t.Unix(), 10) + ".000000"

	api := slack.New(slackAPIToken)

	bulkDelete(api, oldestTimeStamp, "")
}
