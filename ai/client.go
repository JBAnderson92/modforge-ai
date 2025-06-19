package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// Client wraps the OpenAI client
type Client struct {
	client *openai.Client
}

// NewClient creates a new AI client
func NewClient(apiKey string) *Client {
	return &Client{
		client: openai.NewClient(apiKey),
	}
}

// ProcessModRequest represents a request to process a mod
type ProcessModRequest struct {
	Content        string            `json:"content"`
	PromptTemplate string            `json:"prompt_template"`
	GameType       string            `json:"game_type"`
	Variables      map[string]string `json:"variables"`
}

// ProcessModResponse represents the response from processing a mod
type ProcessModResponse struct {
	ProcessedContent string `json:"processed_content"`
	Changelog        string `json:"changelog"`
	TokensUsed       int    `json:"tokens_used"`
}

// ProcessMod processes a mod using AI
func (c *Client) ProcessMod(ctx context.Context, req ProcessModRequest) (*ProcessModResponse, error) {
	// Build the prompt
	prompt := c.buildPrompt(req.PromptTemplate, req.Content, req.Variables)

	// Create the OpenAI request
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o, // Use GPT-4o which is available
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf(`You are an expert game modding assistant specializing in %s mods. 
Your job is to intelligently modify game mod files while preserving their technical structure.

Rules:
1. Always maintain valid JSON/file structure
2. Only modify content that makes sense to change
3. Preserve all technical IDs, keys, and references
4. Provide a brief changelog of what you modified
5. Be conservative - only make improvements that are clearly beneficial

Respond with a JSON object containing:
{
  "processed_content": "the modified content",
  "changelog": "brief summary of changes made"
}`, req.GameType),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   4000,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse the response
	content := resp.Choices[0].Message.Content

	// For now, we'll do simple parsing. In production, we'd parse the JSON response properly
	// This is a placeholder implementation for the MVP
	return &ProcessModResponse{
		ProcessedContent: content,
		Changelog:        "AI-generated modifications applied",
		TokensUsed:       resp.Usage.TotalTokens,
	}, nil
}

// buildPrompt builds the final prompt from template and variables
func (c *Client) buildPrompt(template, content string, variables map[string]string) string {
	prompt := template

	// Replace content placeholder
	prompt = strings.ReplaceAll(prompt, "{content}", content)

	// Replace other variables
	for key, value := range variables {
		placeholder := fmt.Sprintf("{%s}", key)
		prompt = strings.ReplaceAll(prompt, placeholder, value)
	}

	return prompt
}

// ValidateModContent performs basic validation on mod content
func ValidateModContent(content string, gameType string) error {
	switch gameType {
	case "minecraft":
		return validateMinecraftJSON(content)
	case "skyrim":
		return validateSkyrimESP(content)
	case "lua":
		return validateLuaScript(content)
	default:
		return fmt.Errorf("unsupported game type: %s", gameType)
	}
}

// validateMinecraftJSON validates Minecraft JSON mod content
func validateMinecraftJSON(content string) error {
	// Basic JSON validation - in production we'd use proper JSON parsing
	if !strings.Contains(content, "{") || !strings.Contains(content, "}") {
		return fmt.Errorf("invalid JSON structure")
	}
	return nil
}

// validateSkyrimESP validates Skyrim ESP content (placeholder)
func validateSkyrimESP(content string) error {
	// Placeholder validation for Skyrim ESP files
	return nil
}

// validateLuaScript validates Lua script content (placeholder)
func validateLuaScript(content string) error {
	// Placeholder validation for Lua scripts
	return nil
}
