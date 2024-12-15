package service

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"fmt"
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
    reqBody := map[string]interface{}{
        "inputs": context + " " + query,
        "parameters": map[string]interface{}{
            "max_new_tokens": 725,  // Tambahkan parameter ini
            "return_full_text": false,
        },
    }

    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return model.ChatResponse{}, err
    }

    req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct", bytes.NewBuffer(jsonBody))
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
        return model.ChatResponse{}, fmt.Errorf("AI model returned non-OK status: %d", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return model.ChatResponse{}, err
    }

    var chatResp []model.ChatResponse
    err = json.Unmarshal(body, &chatResp)
    if err != nil {
        return model.ChatResponse{}, err
    }

    if len(chatResp) > 0 {
        return chatResp[0], nil
    }

    return model.ChatResponse{}, fmt.Errorf("no response from AI model")
}
