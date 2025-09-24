package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClinetRequest struct {
	Request string `json:"request"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}

type OllamaRespose struct {
	Response string `json:"response"`
}

type Command struct {
	Route   string `json:"route"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func talkToOllama(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"
	reqBody := OllamaRequest{
		Model:  "llama3.1:8b",
		Prompt: prompt,
		Stream: false,
	}
	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	var response OllamaRespose
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", err
	}
	return response.Response, nil
}

func main() {
	// Create a default Gin router with logger + recovery middleware
	r := gin.Default()

	projectMap := map[string]string{
		"project1": "1",
	}

	r.POST("/ai-server", func(c *gin.Context) {
		var req ClinetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		mapJson, _ := json.Marshal(projectMap)
		prompt := fmt.Sprintf(
			`You are a smart command parser.
			The INPUT comes from voice recognition which may contain misheard words or unclear phrases.
			You may receive repeat or similar requests multiple times, just follow the steps to generate the OUTPUT.

 			INPUT: 
			%s

			STEPS:
			- Identify if the user is asking to "broadcast".
			- Get "type": "email" or "sms" from the INPUT.
			- Get "message" (or similar synonyms like "msg") from the INPUT.
			- Find the project id that best matches from this mapping: %s.
			- Default to use email for type param if no match is found.
			- Use captured info to construct the OUTPUT JSON.
			- Respond with only valid JSON.

			DO NOT:
			- Do not add any extra text, notes, reasoning, comments, markdown or explanation.
			- Do not complain or respond in any way other than the OUTPUT.
			- Do not over think.
			
	 		OUTPUT:
			Return only valid JSON in this format without backticks:
			{
			"route": "broadcast/project/{id}",
			"message": "{message}",
			"type": "{email|sms}"
			}`,
			string(mapJson),
			req.Request,
		)

		answer, err := talkToOllama(prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var cmd Command
		println("Ollama response:", answer)
		if err := json.Unmarshal([]byte(answer), &cmd); err != nil {
			c.JSON(422, err)
		}

		c.JSON(200, cmd)
	})

	// Define a simple GET route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start the server on port 8080
	r.Run(":8080") // listen and serve
}
