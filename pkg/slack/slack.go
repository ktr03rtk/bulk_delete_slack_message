package slack

import (
	"fmt"
	"log"
	"strconv"
	"time"

	s "github.com/slack-go/slack"
)

type Client struct {
	*s.Client
}

func NewClient(token string) *Client {
	return &Client{
		s.New(token),
	}
}

func (c *Client) GetChannelNameList(chIDList []string) ([]string, error) {
	channelNameList := make([]string, 0, len(chIDList))

	for _, channelID := range chIDList {
		channel, err := c.GetConversationInfo(channelID, true)
		if err != nil {
			return nil, err
		}

		channelNameList = append(channelNameList, channel.GroupConversation.Name)
	}

	return channelNameList, nil
}

func (c *Client) Delete(timestamp time.Time, channelIDList []string) error {
	for _, channelID := range channelIDList {
		if err := c.delete(timestamp, channelID, ""); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (c *Client) delete(timestamp time.Time, channelID, cursor string) error {
	params := &s.GetConversationHistoryParameters{
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
		ch, t, err := c.DeleteMessage(channelID, msg.Timestamp)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted: channel_ID: %+v, timestamp: %+v\n", ch, t)
		time.Sleep(1200 * time.Millisecond) // chat.delete API is Tier3. Rate limit it 50+ per minute.
	}

	if res.HasMore {
		if err := c.delete(timestamp, channelID, res.ResponseMetaData.NextCursor); err != nil {
			return err
		}
	}

	return nil
}
