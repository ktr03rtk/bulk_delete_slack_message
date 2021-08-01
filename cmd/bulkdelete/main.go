package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

const shortForm = "2006/01/02"

var (
	slackAPIToken string
	channelIDList []string
)

type client struct {
	*slack.Client
}

func getEnv() error {
	s, ok := os.LookupEnv("SLACK_API_TOKEN")
	if !ok {
		return errors.New("env SLACK_API_TOKEN is not found")
	}

	slackAPIToken = s

	c, ok := os.LookupEnv("CHANNEL_ID_LIST")
	if !ok {
		return errors.New("env CHANNEL_ID_LIST is not found")
	}

	channelIDList = strings.Split(c, ",")

	return nil
}

func newClient() *client {
	return &client{
		slack.New(slackAPIToken),
	}
}

func (c *client) getChannelNameList() ([]string, error) {
	channelNameList := make([]string, 0, len(channelIDList))

	for _, channelID := range channelIDList {
		channel, err := c.GetConversationInfo(channelID, true)
		if err != nil {
			return nil, err
		}

		channelNameList = append(channelNameList, channel.GroupConversation.Name)
	}

	return channelNameList, nil
}

func specifyLatestTime() (time.Time, error) {
	fmt.Println("This program delete SLACK messages older than the entered date.")
	fmt.Printf("Enter date in the format like %s:  ", shortForm)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	if input == "" {
		// default: delete messages older than 1 month
		return time.Now().AddDate(0, -1, 0), nil
	}

	t, err := time.Parse(shortForm, input)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func confirm(timestamp string, channelNameList []string) error {
	fmt.Printf("Are you sure you want to delete messages older than %s of Channels %q? (Y/n) >", timestamp, channelNameList)

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

func (c *client) bulkDelete(timestamp time.Time, channelID, cursor string) error {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     100,
		Latest:    strconv.FormatInt(timestamp.Unix(), 10) + ".000000",
		Cursor:    cursor,
	}

	res, err := c.GetConversationHistory(params)
	if err != nil {
		return err
	}

	for _, msg := range res.Messages {
		c, t, err := c.DeleteMessage(channelID, msg.Timestamp)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted: channel_ID: %+v, timestamp: %+v\n", c, t)

		// chat.delete API is Tier3. Rate limit it 50+ per minute.
		time.Sleep(1200 * time.Millisecond)
	}

	if res.HasMore {
		if err := c.bulkDelete(timestamp, channelID, res.ResponseMetaData.NextCursor); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}

	c := newClient()

	channelNameList, err := c.getChannelNameList()
	if err != nil {
		log.Fatal(err)
	}

	latestTimeStamp, err := specifyLatestTime()
	if err != nil {
		log.Fatal(err)
	}

	if err := confirm(latestTimeStamp.Format(shortForm), channelNameList); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for _, channelID := range channelIDList {
		if err := c.bulkDelete(latestTimeStamp, channelID, ""); err != nil {
			log.Fatal(err)
		}
	}
}
