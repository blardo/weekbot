package imagegen

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ImageGenerator handles image generation
type ImageGenerator struct {
	apiKey string
	apiURL string
}

// NewImageGenerator creates a new image generator
func NewImageGenerator(apiKey string) *ImageGenerator {
	if apiKey == "" {
		return nil
	}
	return &ImageGenerator{
		apiKey: apiKey,
		apiURL: "https://api.openai.com/v1/images/generations",
	}
}

// GenerateImage generates an image from a text prompt
// Returns the image data as bytes, or an error
func (ig *ImageGenerator) GenerateImage(prompt string) ([]byte, error) {
	if ig == nil || ig.apiKey == "" {
		return nil, fmt.Errorf("image generator not configured")
	}

	// Enhance the prompt for better results
	enhancedPrompt := fmt.Sprintf("A vibrant, colorful banner image representing: %s. Style: modern, abstract, celebratory, suitable for a Discord server icon", prompt)

	requestBody := map[string]interface{}{
		"model":  "dall-e-3",
		"prompt": enhancedPrompt,
		"n":      1,
		"size":   "1024x1024",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", ig.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ig.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response struct {
		Data []struct {
			URL           string `json:"url"`
			B64JSON       string `json:"b64_json,omitempty"`
			RevisedPrompt string `json:"revised_prompt,omitempty"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no image data in response")
	}

	// If we have base64 data, use it; otherwise fetch from URL
	if response.Data[0].B64JSON != "" {
		imageData, err := base64.StdEncoding.DecodeString(response.Data[0].B64JSON)
		if err != nil {
			return nil, fmt.Errorf("error decoding base64 image: %w", err)
		}
		return imageData, nil
	}

	// Fetch image from URL
	imgResp, err := http.Get(response.Data[0].URL)
	if err != nil {
		return nil, fmt.Errorf("error fetching image from URL: %w", err)
	}
	defer imgResp.Body.Close()

	imageData, err := io.ReadAll(imgResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading image data: %w", err)
	}

	return imageData, nil
}

// GenerateImageForWeek generates an image for a week name
func (ig *ImageGenerator) GenerateImageForWeek(weekName string) ([]byte, error) {
	prompt := fmt.Sprintf("Week name: %s", weekName)
	return ig.GenerateImage(prompt)
}
