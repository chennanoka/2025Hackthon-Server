# Install 

Install AI model management tool Ollama from https://github.com/ollama/ollama

Run "ollama run llama3.1:8b" from cmd (ref:https://ollama.com/library/llama3.1:8b)

Install golang 

Run "go run main.go" at the root of the folder to start server.

Then, you can POST http://localhost:8080/ai-server with JSON e.g:
{
    "request": "broadcast to project1 with message It's a good day nice via SMS"
}

to test the server.

It should return e.g 

{
    "route": "broadcast/project/1",
    "extra": "It's a good day nice",
    "type": "sms"
}