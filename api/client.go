package api

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type ErrorRequest struct {
    ErrorDetails string `json:"error_details"`
    Context      string `json:"context"`
}

type ErrorResponse struct {
    Resolutions []string `json:"resolutions"`
    Info        string   `json:"info"`
}

const baseURL = "http://localhost:8080/api"

func SendErrorRequest(request ErrorRequest) (*ErrorResponse, error) {
    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(baseURL+"/errors", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var response ErrorResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, err
    }

    return &response, nil
}

func HealthCheck() (string, error) {
    resp, err := http.Get(baseURL + "/health")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var status string
    if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
        return "", err
    }

    return status, nil
}