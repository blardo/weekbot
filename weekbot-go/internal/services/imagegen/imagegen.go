package imagegen

import (
	"context"
	"fmt"

	"weekbot-go/internal/logger"

	"google.golang.org/genai"
)

// ImageGenerator handles image generation using Google Gemini
type ImageGenerator struct {
	client *genai.Client
	model  string
}

// NewImageGenerator creates a new image generator using Gemini API
func NewImageGenerator(apiKey string) *ImageGenerator {
	if apiKey == "" {
		logger.Debug("ImageGenerator: API key not provided, returning nil")
		return nil
	}

	ctx := context.Background()
	
	logger.Debug("ImageGenerator: Creating Gemini API client", "model", "gemini-2.5-flash-image")
	// Create client with API key
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.Error("ImageGenerator: Failed to create Gemini client", "error", err)
		// Return nil if client creation fails - will be handled gracefully
		return nil
	}

	logger.Info("ImageGenerator: Successfully created Gemini client")
	return &ImageGenerator{
		client: client,
		model:  "gemini-2.5-flash-image", // Fast image generation model
	}
}

// GenerateImage generates an image from a text prompt using Gemini
// Returns the image data as bytes, or an error
func (ig *ImageGenerator) GenerateImage(prompt string) ([]byte, error) {
	if ig == nil || ig.client == nil {
		logger.Error("ImageGenerator: GenerateImage called but generator not configured")
		return nil, fmt.Errorf("image generator not configured")
	}

	// Enhance the prompt for better results
	enhancedPrompt := fmt.Sprintf("Create a vibrant, colorful banner image representing: %s. Style: modern, abstract, celebratory, suitable for a Discord server banner (16:9 aspect ratio)", prompt)

	logger.Info("ImageGenerator: Starting image generation", "model", ig.model, "prompt", prompt, "enhanced_prompt", enhancedPrompt)
	ctx := context.Background()
	
	// Generate images using Gemini image generation model
	// Use GenerateContent - the model will automatically return images when using image generation models
	result, err := ig.client.Models.GenerateContent(ctx, ig.model, genai.Text(enhancedPrompt), nil)
	if err != nil {
		logger.Error("ImageGenerator: API call failed", "error", err, "model", ig.model, "prompt", prompt)
		return nil, fmt.Errorf("error generating image: %w", err)
	}

	logger.Debug("ImageGenerator: API call succeeded", "candidates_count", len(result.Candidates))

	if len(result.Candidates) == 0 {
		logger.Error("ImageGenerator: No candidates in API response")
		return nil, fmt.Errorf("no candidates in response")
	}

	// Extract image data from response parts
	candidate := result.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		logger.Error("ImageGenerator: No parts in candidate content")
		return nil, fmt.Errorf("no parts in response")
	}

	// Look for image data in the parts
	for _, part := range candidate.Content.Parts {
		if part.InlineData != nil && len(part.InlineData.Data) > 0 {
			logger.Info("ImageGenerator: Image generated successfully", "size_bytes", len(part.InlineData.Data), "mime_type", part.InlineData.MIMEType)
			return part.InlineData.Data, nil
		}
	}

	logger.Error("ImageGenerator: No image data found in response parts")
	return nil, fmt.Errorf("no image data found in response")
}

// GenerateImageForWeek generates an image for a week name
func (ig *ImageGenerator) GenerateImageForWeek(weekName string) ([]byte, error) {
	logger.Info("ImageGenerator: Generating image for week", "week_name", weekName)
	prompt := fmt.Sprintf("Week name: %s", weekName)
	return ig.GenerateImage(prompt)
}
