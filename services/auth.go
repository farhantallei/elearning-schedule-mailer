package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AuthPayload struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AuthResponseResult struct {
	Token         string `json:"token"`
	Username      string `json:"username"`
	CryptUsername string `json:"crypt_username"`
}

type AuthResponseBody struct {
	Result AuthResponseResult `json:"result"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	NIM      string `json:"nim"`
	CryptNIM string `json:"crypt_nim"`
}

func auth(baseURL string, clientId string, clientSecret string) (AuthResponse, error) {
	var authResponse AuthResponse

	username := os.Getenv("E_LEARNING_USERNAME")
	password := os.Getenv("E_LEARNING_PASSWORD")

	if username == "" || password == "" {
		return authResponse, fmt.Errorf("missing environment variables")
	}

	url := fmt.Sprintf("%s/user/login", baseURL)

	payload := AuthPayload{
		Username:     username,
		Password:     password,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return authResponse, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return authResponse, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "api-elearning.utb-univ.ac.id")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return authResponse, fmt.Errorf("error sending request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return authResponse, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return authResponse, fmt.Errorf("error reading response body: %w", err)
	}

	var authResponseBody AuthResponseBody

	err = json.Unmarshal(responseBody, &authResponseBody)
	if err != nil {
		return authResponse, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	authResponse.Token = authResponseBody.Result.Token
	authResponse.NIM = authResponseBody.Result.Username
	authResponse.CryptNIM = authResponseBody.Result.CryptUsername

	return authResponse, nil
}
