package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type MailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Html    string `json:"html"`
}

type MailResponseError struct {
	Error string `json:"error"`
}

func SendMail(subject string, body string) (*string, error) {
	url := "https://react.email/api/send/test"

	myEmail := os.Getenv("MY_EMAIL")
	if myEmail == "" {
		return nil, fmt.Errorf("missing environment variables")
	}

	payload := MailPayload{
		To:      myEmail,
		Subject: subject,
		Html:    body,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending mail: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode == 429 {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}
		var mailResponseError MailResponseError
		err = json.Unmarshal(responseBody, &mailResponseError)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling response body: %w", err)
		}
		return nil, fmt.Errorf("error sending mail: %s", mailResponseError.Error)
	}

	mailResponse := "Mail sent successfully"

	return &mailResponse, nil
}
