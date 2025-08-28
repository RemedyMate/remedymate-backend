package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Example client for testing the conversation service
func main() {
	baseURL := "http://localhost:8080/api/v1"

	fmt.Println("🚀 RemedyMate Conversation Service Example")
	fmt.Println("==========================================")

	// Step 1: Start a conversation (no authentication required)
	fmt.Println("\n1. Starting conversation...")
	startReq := map[string]interface{}{
		"symptom":  "headache",
		"language": "en",
		// Note: user_id is optional for unauthenticated users
	}

	startResp, err := makeRequest("POST", baseURL+"/conversation/start", startReq)
	if err != nil {
		fmt.Printf("❌ Failed to start conversation: %v\n", err)
		return
	}

	conversationID := startResp["conversation_id"].(string)
	fmt.Printf("✅ Conversation started: %s\n", conversationID)
	fmt.Printf("📝 First question: %s\n", startResp["question"].(map[string]interface{})["text"])

	// Step 2: Submit answers
	answers := []string{
		"3 days",           // Duration
		"Front of my head", // Location
		"Moderate",         // Severity
		"None",             // Medical history
		"Stress",           // Triggers
	}

	for i, answer := range answers {
		fmt.Printf("\n%d. Submitting answer: %s\n", i+2, answer)

		answerReq := map[string]interface{}{
			"conversation_id": conversationID,
			"answer":          answer,
		}

		answerResp, err := makeRequest("POST", baseURL+"/conversation/answer", answerReq)
		if err != nil {
			fmt.Printf("❌ Failed to submit answer: %v\n", err)
			return
		}

		if answerResp["is_complete"].(bool) {
			fmt.Println("✅ All questions completed!")
			break
		}

		question := answerResp["question"].(map[string]interface{})
		fmt.Printf("📝 Next question: %s\n", question["text"])
	}

	// Step 3: Get the final report
	fmt.Println("\n3. Getting final health report...")
	reportResp, err := makeRequest("GET", baseURL+"/conversation/"+conversationID+"/report", nil)
	if err != nil {
		fmt.Printf("❌ Failed to get report: %v\n", err)
		return
	}

	fmt.Println("📊 Final Health Report:")
	fmt.Println("=======================")
	report := reportResp["report"].(map[string]interface{})
	fmt.Printf("Symptom: %s\n", report["symptom"])
	fmt.Printf("Duration: %s\n", report["duration"])
	fmt.Printf("Location: %s\n", report["location"])
	fmt.Printf("Severity: %s\n", report["severity"])
	fmt.Printf("Urgency Level: %s\n", report["urgency_level"])

	if conditions, ok := report["possible_conditions"].([]interface{}); ok {
		fmt.Printf("Possible Conditions: %v\n", conditions)
	}

	if recommendations, ok := report["recommendations"].([]interface{}); ok {
		fmt.Printf("Recommendations: %v\n", recommendations)
	}

	fmt.Println("\n✅ Conversation flow completed successfully!")
}

// Helper function to make HTTP requests
func makeRequest(method, url string, body interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
