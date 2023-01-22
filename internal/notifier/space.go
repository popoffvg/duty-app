package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Payload struct {
	Recipient string  `json:"recipient"`
	Content   Content `json:"content"`
}

type Content struct {
	ClassName string `json:"className"`
	Text      string `json:"text"`
}

type SpaceClient struct {
	endpoint string
	token    string
}

func NewSpaceClient(endpoint, token string) (*SpaceClient, error) {
	if endpoint == "" || token == "" {
		return nil, fmt.Errorf("space client: endpoint or token is empty")
	}

	return &SpaceClient{
		endpoint: endpoint,
		token:    token,
	}, nil
}

func (sc *SpaceClient) SendNotification(channelName, text string) error {
	data := Payload{
		Recipient: "channel:name:" + channelName,
		Content: Content{
			ClassName: "ChatMessage.Text",
			Text:      text,
		},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json marshall error: %w", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", sc.endpoint, body)
	if err != nil {
		return fmt.Errorf("http request error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+sc.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http client error: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
