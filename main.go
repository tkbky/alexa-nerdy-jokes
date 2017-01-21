package main

import (
	"log"
	"os"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

var applications = map[string]interface{}{
	"/echo/helloworld": alexa.EchoApplication{
		AppID:    os.Getenv("ALEXA_SKILL_APP_ID"),
		OnIntent: HelloWorldHandler,
		OnLaunch: HelloWorldHandler,
	},
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	alexa.Run(applications, port)
}

// HelloWorldHandler says "Hello"
func HelloWorldHandler(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	echoResp.OutputSpeech("Hi, What can I help?").Card("Greeting", "This is a greeting")
}
