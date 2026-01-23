package imagegen

import (
	"context"
	"fmt"

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
		return nil
	}

	ctx := context.Background()
	
	// Create client with API key
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		// Return nil if client creation fails - will be handled gracefully
		return nil
	}

	return &ImageGenerator{
		client: client,
		model:  "gemini-2.5-flash-image", // Fast image generation model
	}
}

// GenerateImage generates an image from a text prompt using Gemini
// Returns the image data as bytes, or an error
func (ig *ImageGenerator) GenerateImage(prompt string) ([]byte, error) {
	if ig == nil || ig.client == nil {
		return nil, fmt.Errorf("image generator not configured")
	}

	// Enhance the prompt for better results
	enhancedPrompt := fmt.Sprintf("Create a vibrant, colorful banner image representing: %s. Style: modern, abstract, celebratory, suitable for a Discord server icon (1024x1024 pixels)", prompt)

	ctx := context.Background()
	
	// Generate images using Gemini image generation model
	result, err := ig.client.Models.GenerateImages(ctx, ig.model, enhancedPrompt, &genai.GenerateImagesConfig{
		NumberOfImages: 1,
		AspectRatio:    "1:1",
	})
	if err != nil {
		return nil, fmt.Errorf("error generating image: %w", err)
	}

	if len(result.GeneratedImages) == 0 {
		return nil, fmt.Errorf("no images in response")
	}

	// Extract image data from response
	generatedImage := result.GeneratedImages[0]
	if generatedImage.Image == nil {
		return nil, fmt.Errorf("no image data in generated image")
	}

	// Prefer ImageBytes if available, otherwise fetch from GCSURI
	if len(generatedImage.Image.ImageBytes) > 0 {
		return generatedImage.Image.ImageBytes, nil
	}

	if generatedImage.Image.GCSURI != "" {
		return nil, fmt.Errorf("image stored at GCS URI, not supported: %s", generatedImage.Image.GCSURI)
	}

	return nil, fmt.Errorf("no image data found in response")
}

// GenerateImageForWeek generates an image for a week name
func (ig *ImageGenerator) GenerateImageForWeek(weekName string) ([]byte, error) {
	prompt := fmt.Sprintf("Week name: %s", weekName)
	return ig.GenerateImage(prompt)
}
