package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GroqRequest struct {
	Model    string        `json:"model"`
	Messages []GroqMessage `json:"messages"`
}

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CallGroq(prompt string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	requestBody := GroqRequest{
		Model: "llama3-8b-8192",
		Messages: []GroqMessage{
			{Role: "system", Content: "You are a career coach that gives humorous yet useful feedback."},
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ os.Getenv("GROQ_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Groq API error: " + string(body))
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", err
	}

	return groqResp.Choices[0].Message.Content, nil
}

func GeneratePrompt(resumeText string, mode string) string {
	if mode == "roast" {
		return "You are a brutally honest, sarcastic, and funny AI career coach. You're reviewing a resume and giving roast-style feedback—like a stand-up comedian who knows hiring. Be witty, sharp, and pull no punches. Mix humor with real insights. The goal is to mock obvious resume clichés, vague fluff, and poor formatting, but still give helpful suggestions beneath the jokes. Avoid any sugarcoating. This resume is yours: \n\n" + resumeText
	}
	return "Give constructive career feedback on this resume and add a little joke here and there but dont be to rough :\n\n" + resumeText
}
