package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ktr03rtk/bulk_delete_slack_message/pkg/slack"
	"github.com/pkg/errors"
)

const shortForm = "2006/01/02"

var (
	slackAPIToken string
	channelList   []string
)

func getEnv() error {
	s, ok := os.LookupEnv("SLACK_API_TOKEN")
	if !ok {
		return errors.New("env SLACK_API_TOKEN is not found")
	}

	slackAPIToken = s

	c, ok := os.LookupEnv("CHANNEL_LIST")
	if !ok {
		return errors.New("env CHANNEL_LIST is not found")
	}

	channelList = strings.Split(c, ",")

	return nil
}

func getTargetTime() (*time.Time, error) {
	fmt.Println("This program delete SLACK messages older than the date you enter.")
	fmt.Printf("Enter date in the format like %s:  ", shortForm)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// default: delete messages older than 1 month
	if input == "" {
		t := time.Now().AddDate(0, -1, 0)

		return &t, nil
	}

	t, err := time.Parse(shortForm, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse timestamp")
	}

	return &t, nil
}

func confirm(timestamp string, channelList []string) error {
	fmt.Printf("Are you sure you want to delete messages of Channels %q older than %s? (Y/n) >", channelList, timestamp)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	switch strings.ToLower(scanner.Text()) {
	case "y", "yes":
		fmt.Println("start.")
	default:
		return errors.New("aborting the process")
	}

	return nil
}

func main() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}

	c := slack.NewClient(slackAPIToken)

	channelMap, err := c.GetChannelMap(channelList)
	if err != nil {
		log.Fatal(err)
	}

	targetTimestamp, err := getTargetTime()
	if err != nil {
		log.Fatal(err)
	}

	if err := confirm(targetTimestamp.Format(shortForm), channelList); err != nil {
		log.Fatal(err)
	}

	if err := c.Delete(*targetTimestamp, channelMap); err != nil {
		log.Fatal(err)
	}
}
