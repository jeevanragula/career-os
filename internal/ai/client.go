package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

type Client struct { BaseURL, APIKey, Model string; HTTP *http.Client }

func (c Client) Complete(ctx context.Context, system, user string) (string, error) {
    if c.BaseURL == "" || c.Model == "" {
        return "", fmt.Errorf("AI_BASE_URL and AI_MODEL are required")
    }

    // Gemini 3.8 Flash does not support the legacy temperature parameter.
    body := map[string]any{
        "model": c.Model,
        "messages": []any{
            map[string]string{"role": "system", "content": system},
            map[string]string{"role": "user", "content": user},
        },
    }
    raw, _ := json.Marshal(body)

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        strings.TrimRight(c.BaseURL, "/")+"/chat/completions",
        bytes.NewReader(raw),
    )
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")
    if c.APIKey != "" {
        req.Header.Set("Authorization", "Bearer "+c.APIKey)
    }

    client := c.HTTP
    if client == nil {
        client = &http.Client{Timeout: 120 * time.Second}
    }

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 16*1024))
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        detail := strings.TrimSpace(string(responseBody))
        if detail == "" {
            return "", fmt.Errorf("AI HTTP %d", resp.StatusCode)
        }
        return "", fmt.Errorf("AI HTTP %d: %s", resp.StatusCode, detail)
    }

    var out map[string]any
    if err := json.Unmarshal(responseBody, &out); err != nil {
        return "", err
    }

    choices, _ := out["choices"].([]any)
    if len(choices) == 0 {
        return "", fmt.Errorf("AI returned no choices")
    }

    msg, _ := choices[0].(map[string]any)
    message, _ := msg["message"].(map[string]any)
    content, _ := message["content"].(string)
    if content == "" {
        return "", fmt.Errorf("AI returned empty content")
    }

    return content, nil
}
