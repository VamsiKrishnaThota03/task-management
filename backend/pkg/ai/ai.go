package ai

import (
    "context"
    "fmt"
    "log"
    "strings"
    "google.golang.org/api/option"
    "github.com/google/generative-ai-go/genai"
)

func GetTaskSuggestions(prompt string, apiKey string) (string, error) {
    ctx := context.Background()

    // Create client
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        log.Printf("Failed to create client: %v", err)
        return "", fmt.Errorf("failed to create Gemini client: %v", err)
    }
    defer client.Close()

    // Create model
    model := client.GenerativeModel("gemini-pro")

    // Configure the generation
    model.SetTemperature(0.7)
    model.SetTopK(40)
    model.SetTopP(0.95)
    model.SetMaxOutputTokens(150)

    // Create a more specific prompt
    prompt = `Please provide brief and concise suggestions for the following task description. 
    Format your response in simple text with:
    1. 2-3 suggestions for improving the task description
    2. 3-4 specific subtasks
    
    Task description: ` + prompt

    // Send the prompt directly without chat session
    resp, err := model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        log.Printf("Failed to get response: %v", err)
        return "", fmt.Errorf("failed to get suggestions: %v", err)
    }

    if len(resp.Candidates) == 0 {
        return "", fmt.Errorf("no response received from Gemini")
    }

    // Extract text from response
    var responseText string
    for _, part := range resp.Candidates[0].Content.Parts {
        if textPart, ok := part.(genai.Text); ok {
            responseText += string(textPart)
        }
    }

    if responseText == "" {
        return "", fmt.Errorf("no text content in Gemini response")
    }

    // Clean up the response
    responseText = strings.ReplaceAll(responseText, "**", "")
    responseText = strings.ReplaceAll(responseText, "*", "•")
    
    // Format the response
    lines := strings.Split(responseText, "\n")
    var formattedLines []string
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" {
            formattedLines = append(formattedLines, line)
        }
    }
    
    return strings.Join(formattedLines, "\n"), nil
}
