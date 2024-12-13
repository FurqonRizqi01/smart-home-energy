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
    if len(table) == 0 {
        return "", fmt.Errorf("table is empty")
    }

    reqBody := model.AIRequest{
        Inputs: model.Inputs{
            Table: table,
            Query: query,
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
        return "", fmt.Errorf("AI model returned non-OK status: %d", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    fmt.Println("Raw Response:", string(body))

    var tapasResp model.TapasResponse
    err = json.Unmarshal(body, &tapasResp)
    if err != nil {
        return "", fmt.Errorf("failed to unmarshal response: %v", err)
    }

    if len(tapasResp.Cells) > 0 {
        return tapasResp.Cells[0], nil
    }

    return "", fmt.Errorf("no response from AI model")
}


func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
    reqBody := map[string]string{
        "inputs": context + " " + query,
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
