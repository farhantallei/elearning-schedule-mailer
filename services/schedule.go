package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ScheduleResponseResult struct {
	CourseName   string  `json:"course_name"`
	LecturerName string  `json:"lecturer_name"`
	CourseTopic  string  `json:"course_topic"`
	Noted        *string `json:"noted"`
	LinkMedia    *string `json:"link_media"`
}

type ScheduleResponseBody struct {
	Result []ScheduleResponseResult `json:"result"`
}

func FetchSchedule(courseName string) (ScheduleResponseResult, error) {
	var scheduleResponse ScheduleResponseResult

	baseURL := os.Getenv("E_LEARNING_API_URL")
	clientId := os.Getenv("E_LEARNING_CLIENT_ID")
	clientSecret := os.Getenv("E_LEARNING_CLIENT_SECRET")

	if baseURL == "" || clientId == "" || clientSecret == "" {
		return scheduleResponse, fmt.Errorf("missing environment variables")
	}

	authResponse, err := auth(baseURL, clientId, clientSecret)
	if err != nil {
		return scheduleResponse, fmt.Errorf("error authenticating: %w", err)
	}

	dayNames := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}

	today := time.Now()

	params := url.Values{}
	params.Add("client_id", clientId)
	params.Add("client_secret", clientSecret)
	params.Add("token", authResponse.Token)
	params.Add("nim", authResponse.CryptNIM)
	params.Add("day_name", dayNames[today.Weekday()])
	params.Add("date_now", today.Format("2006-01-02"))

	fullURL := fmt.Sprintf("%s?%s", fmt.Sprintf("%s/eLearning/courseTopic-student/session", baseURL), params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return scheduleResponse, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Host", "api-elearning.utb-univ.ac.id")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return scheduleResponse, fmt.Errorf("error sending request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing response body: ", err)
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return scheduleResponse, fmt.Errorf("error reading response body: %v", err)
	}

	var scheduleResponseBody ScheduleResponseBody
	err = json.Unmarshal(responseBody, &scheduleResponseBody)
	if err != nil {
		return scheduleResponse, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	for _, result := range scheduleResponseBody.Result {
		if result.CourseName == courseName {
			return result, nil
		}
	}

	return scheduleResponse, fmt.Errorf("course not found")
}
