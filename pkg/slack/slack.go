package slack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	s "github.com/slack-go/slack"
)

type (
	ChannelName string
	ChannelID   string
	ChannelMap  map[ChannelName]ChannelID
)

type Client struct {
	*s.Client
}

func NewClient(token string) *Client {
	return &Client{
		s.New(token),
	}
}

func (c *Client) GetChannelMap(targetChList []string) (ChannelMap, error) {
	result := make(ChannelMap)

	chList, _, err := c.GetConversations(&s.GetConversationsParameters{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get conversations")
	}

	for _, targetCh := range targetChList {
		chID, err := findChannelID(chList, targetCh)
		if err != nil {
			return nil, err
		}

		result[ChannelName(targetCh)] = chID
	}

	return result, nil
}

func findChannelID(chList []s.Channel, targetCh string) (ChannelID, error) {
	for _, ch := range chList {
		if ch.GroupConversation.Name == targetCh {
			return ChannelID(ch.GroupConversation.Conversation.ID), nil
		}
	}

	return "", errors.Errorf("failed to find channel: \"%s\"", targetCh)
}

func (c *Client) Delete(timestamp time.Time, channelMap ChannelMap) error {
	for channelName, channelID := range channelMap {
		if err := c.delete(timestamp, channelName, channelID, ""); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) delete(timestamp time.Time, channelName ChannelName, channelID ChannelID, nextCursor string) error {
	params := &s.GetConversationHistoryParameters{
		ChannelID: string(channelID),
		Limit:     100,
		Latest:    strconv.FormatInt(timestamp.Unix(), 10) + ".000000",
		Cursor:    nextCursor,
	}

	res, err := c.GetConversationHistory(params)
	if err != nil {
		return errors.Wrapf(err, "failed to get conversation history")
	}

	for _, msg := range res.Messages {
		_, t, err := c.DeleteMessage(string(channelID), msg.Timestamp)
		if err != nil {
			return errors.Wrapf(err, "failed to delete message")
		}

		unixTime, err := strconv.ParseInt(strings.Split(t, ".")[0], 10, 64)
		if err != nil {
			return errors.Wrapf(err, "failed to parse timestamp")
		}

		fmt.Printf("Deleted. channel: %+v, timestamp: %+v\n", channelName, time.Unix(unixTime, 0).Local())
		time.Sleep(1200 * time.Millisecond) // chat.delete API is Tier3. Rate limit is 50+ per minute.
	}

	if res.HasMore {
		if err := c.delete(timestamp, channelName, channelID, res.ResponseMetaData.NextCursor); err != nil {
			return err
		}
	}

	return nil
}
