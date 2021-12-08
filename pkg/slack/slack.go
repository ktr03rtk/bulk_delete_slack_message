package slack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	s "github.com/slack-go/slack"
)

type ChannelIDMap map[string]string

type Client struct {
	*s.Client
}

func NewClient(token string) *Client {
	return &Client{
		s.New(token),
	}
}

func (c *Client) GetChannelIDMap(targetChList []string) (ChannelIDMap, error) {
	result := make(ChannelIDMap)

	chList, _, err := c.GetConversations(&s.GetConversationsParameters{})
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	for _, targetCh := range targetChList {
		chID, err := findChannelID(chList, targetCh)
		if err != nil {
			return nil, err
		}

		result[targetCh] = chID
	}

	return result, nil
}

func findChannelID(chList []s.Channel, targetCh string) (string, error) {
	for _, ch := range chList {
		if ch.GroupConversation.Name == targetCh {
			return ch.GroupConversation.Conversation.ID, nil
		}
	}

	return "", fmt.Errorf("failed to find channel: \"%s\"", targetCh)
}

func (c *Client) Delete(timestamp time.Time, channelIDMap ChannelIDMap) error {
	for channelName, channelID := range channelIDMap {
		if err := c.delete(timestamp, channelName, channelID, ""); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) delete(timestamp time.Time, channelName, channelID, nextCursor string) error {
	params := &s.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     100,
		Latest:    strconv.FormatInt(timestamp.Unix(), 10) + ".000000",
		Cursor:    nextCursor,
	}

	res, err := c.GetConversationHistory(params)
	if err != nil {
		return fmt.Errorf("failed to get conversation history: %w", err)
	}

	for _, msg := range res.Messages {
		_, t, err := c.DeleteMessage(channelID, msg.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}

		unixTime, err := strconv.ParseInt(strings.Split(t, ".")[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %w", err)
		}

		fmt.Printf("Deleted. channel: %+v, timestamp: %+v\n", channelName, time.Unix(unixTime, 0).Local())
		time.Sleep(1200 * time.Millisecond) // chat.delete API is Tier3. Rate limit it 50+ per minute.
	}

	if res.HasMore {
		if err := c.delete(timestamp, channelName, channelID, res.ResponseMetaData.NextCursor); err != nil {
			return err
		}
	}

	return nil
}
