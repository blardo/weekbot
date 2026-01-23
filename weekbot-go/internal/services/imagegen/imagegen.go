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
	result, err := ig.client.Models.GenerateImages(ctx, ig.model, enhancedPrompt, &genai.GenerateImagesConfig{
		NumberOfImages: 1,
		AspectRatio:    "16:9", // Banner aspect ratio
	})
	if err != nil {
		logger.Error("ImageGenerator: API call failed", "error", err, "model", ig.model, "prompt", prompt)
		return nil, fmt.Errorf("error generating image: %w", err)
	}

	logger.Debug("ImageGenerator: API call succeeded", "images_count", len(result.GeneratedImages))

	if len(result.GeneratedImages) == 0 {
		logger.Error("ImageGenerator: No images in API response")
		return nil, fmt.Errorf("no images in response")
	}

	// Extract image data from response
	generatedImage := result.GeneratedImages[0]
	if generatedImage.Image == nil {
		logger.Error("ImageGenerator: Generated image has nil Image field")
		return nil, fmt.Errorf("no image data in generated image")
	}

	// Prefer ImageBytes if available, otherwise fetch from GCSURI
	if len(generatedImage.Image.ImageBytes) > 0 {
		logger.Info("ImageGenerator: Image generated successfully", "size_bytes", len(generatedImage.Image.ImageBytes), "mime_type", generatedImage.Image.MIMEType)
		return generatedImage.Image.ImageBytes, nil
	}

	if generatedImage.Image.GCSURI != "" {
		logger.Error("ImageGenerator: Image stored at GCS URI (not supported)", "gcs_uri", generatedImage.Image.GCSURI)
		return nil, fmt.Errorf("image stored at GCS URI, not supported: %s", generatedImage.Image.GCSURI)
	}

	logger.Error("ImageGenerator: No image data found in response")
	return nil, fmt.Errorf("no image data found in response")
}

// GenerateImageForWeek generates an image for a week name
func (ig *ImageGenerator) GenerateImageForWeek(weekName string) ([]byte, error) {
	logger.Info("ImageGenerator: Generating image for week", "week_name", weekName)
	prompt := fmt.Sprintf("Week name: %s", weekName)
	return ig.GenerateImage(prompt)
}
