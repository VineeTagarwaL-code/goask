package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type AIRequest struct {
	Contents []Content `json:"contents"`
}

type Candidate struct {
	Content struct {
		Parts []Part `json:"parts"`
		Role  string `json:"role"`
	} `json:"content"`
	FinishReason string `json:"finishReason"`
}

type AIResponse struct {
	Candidates []Candidate `json:"candidates"`
}

func getAIResponse(userPrompt string) (string, error) {

	APIKey := os.Getenv("GEMINI_API_KEY")
	APIEndpoint := os.Getenv("GEMINI_API_ENDPOINT")

	if APIKey == "" || APIEndpoint == "" {
		return "", fmt.Errorf("API key or endpoint not set")
	}
	systemPrompt := `You are a Linux command assistant. Your task is to provide the exact Linux command that solves the user's request. Follow these rules strictly:

1. If the user's request is about a Linux command or task, respond with ONLY the specific Linux command to accomplish it. Do not include any explanations or additional text.
2. If the user's request is not related to a Linux command or task, respond with exactly "Sorry, I can't help with that. I only provide Linux commands."
3. If multiple commands are needed, combine them into a single line using && or ; as appropriate.
4. Do not include any markdown formatting, quotation marks, or code blocks in your response.

User request: `

	fullPrompt := systemPrompt + userPrompt

	request := AIRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: fullPrompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s?key=%s", APIEndpoint, APIKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var aiResp AIResponse
	err = json.Unmarshal(body, &aiResp)
	if err != nil {
		return "", err
	}

	if len(aiResp.Candidates) > 0 && len(aiResp.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(aiResp.Candidates[0].Content.Parts[0].Text), nil
	}

	return "", fmt.Errorf("no response from AI")
}

func runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	err := godotenv.Load()
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s how to do something in linux\n", filepath.Base(os.Args[0]))
		return
	}

	userPrompt := strings.Join(os.Args[2:], " ")
	fmt.Printf("Asking AI for Linux command: %s\n", userPrompt)

	aiResponse, err := getAIResponse(userPrompt)
	if err != nil {
		fmt.Printf("Error getting AI response: %v\n", err)
		return
	}

	if aiResponse == "Sorry, I can't help with that. I only provide Linux commands." {
		fmt.Println(aiResponse)
		return
	}

	fmt.Printf("AI suggests command: %s\n", aiResponse)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to run this command? (y/n): ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(strings.ToLower(userInput))

	if userInput == "y" {
		fmt.Println("Executing command...")
		err := runCommand(aiResponse)
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		} else {
			fmt.Println("Command executed successfully.")
		}
	} else {
		fmt.Println("Command not executed.")
	}
}
