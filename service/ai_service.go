package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	// Validasi input
	if len(table) == 0 {
		return "", fmt.Errorf("table is empty")
	}

	processedTable := make([][]string, 0)

	headers := make([]string, 0)
	for header := range table {
		headers = append(headers, header)
	}
	processedTable = append(processedTable, headers)

	// Tambahkan data
	rowCount := len(table[headers[0]])
	for i := 0; i < rowCount; i++ {
		row := make([]string, len(headers))
		for j, header := range headers {
			row[j] = table[header][i]
		}
		processedTable = append(processedTable, row)
	}

	reqBody := map[string]interface{}{
		"inputs": map[string]interface{}{
			"table": processedTable,
			"query": query,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("AI model returned non-OK status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Raw Tapas Response:", string(body))

	var tapasResp map[string]interface{}
	err = json.Unmarshal(body, &tapasResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	var answer string
	switch v := tapasResp["answer"].(type) {
	case string:
		answer = v
	case []interface{}:
		if len(v) > 0 {
			answer = fmt.Sprintf("%v", v[0])
		}
	}

	if answer == "" {
		if cells, ok := tapasResp["cells"].([]interface{}); ok && len(cells) > 0 {
			answer = fmt.Sprintf("%v", cells[0])
		}
	}

	if answer == "" {
		answer = "No specific answer could be extracted from the response."
	}

	return answer, nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	messages := make([]map[string]string, 0, 2)
	if context != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": context,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": query,
	})

	reqBody := map[string]interface{}{
		"model":      "Qwen/Qwen3-4B-Instruct-2507:cheapest",
		"messages":   messages,
		"max_tokens": 512,
		"stream":     false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return model.ChatResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://router.huggingface.co/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return model.ChatResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return model.ChatResponse{}, fmt.Errorf("AI model returned non-OK status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.ChatResponse{}, err
	}

	var routerResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &routerResp); err == nil &&
		len(routerResp.Choices) > 0 && routerResp.Choices[0].Message.Content != "" {
		return model.ChatResponse{GeneratedText: routerResp.Choices[0].Message.Content}, nil
	}

	// Keep accepting the legacy response shape so existing tests/mocks remain valid.
	var legacyResp []model.ChatResponse
	if err := json.Unmarshal(body, &legacyResp); err == nil &&
		len(legacyResp) > 0 && legacyResp[0].GeneratedText != "" {
		return legacyResp[0], nil
	}

	return model.ChatResponse{}, fmt.Errorf("AI model returned an invalid response: %s", string(body))
}
