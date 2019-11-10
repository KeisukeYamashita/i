package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Message ...
type Message struct {
	Text        string        `json:"text"`
	Attachments *[]Attachment `json:"attachments"`
}

// Attachment ...
type Attachment struct {
	Fields []Field `json:"fields"`
	Title  string  `json:"title"`
}

// Field ...
type Field struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
}

// Client ...
type Client interface {
	PostMessage() error
}

var _ Client = (*client)(nil)

type client struct {
	HookURL string
}

// NewClient ...
func NewClient(url string) Client {
	return &client{
		HookURL: url,
	}
}

func (c *client) PostMessage() error {
	client := http.DefaultClient
	msg := &Message{
		Text: "hello",
	}
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		c.HookURL,
		bytes.NewBuffer(json),
	)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("response from slack api was not 200 got %d", resp.StatusCode)
	}

	return nil
}
