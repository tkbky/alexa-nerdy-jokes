package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

var nerdyJokes = []string{
	"What do you get when you put root beer in a square class? Beer",
	"C, E Flat, and G walk into a bar. The bartender says, \"Sorry, no minors.\"",
	"The past, the present, and the future walked into a bar. It was tense.",
	"What's another name for santa's elves? It's subordinate clauses.",
	"Helium walks into a bar and orders a beer. The bartender says, \"Sorry, We don't serve noble gases here.\" He doesn't react",
	"A photon checks into a hotel and the bellhop asks him if he has any luggage. The photon replies, \"No, I'm travelling light.\"",
	"Why can't you trust atom? because they make up everything?",
	"The first rule of tautology club is the first rule of tautology club.",
	"A biologist, a chemist, and a statistician are out for hunting. The biologist shoots at a deer and misses 5 feet to the left. The chemist shoots and misses 5 feet to the right. The statistician yells, \"We got them!\"",
}

var applications = map[string]interface{}{
	"/echo/nerdyjokes": alexa.EchoApplication{
		AppID:   os.Getenv("ALEXA_SKILL_APP_ID"),
		Handler: NerdyJokesHandler,
	},
}

var db *sql.DB

func main() {
	var err error

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS jokes (id SERIAL NOT NULL PRIMARY KEY, content text NOT NULL, created_at timestamp DEFAULT CURRENT_TIMESTAMP, updated_at timestamp DEFAULT CURRENT_TIMESTAMP)")

	if err != nil {
		log.Fatal(err)
	}

	var count int

	err = db.QueryRow("SELECT COUNT(id) FROM jokes").Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	if count <= 0 {
		for _, joke := range nerdyJokes {
			seedDB(db, joke)
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	alexa.Run(applications, port)
}

func seedDB(db *sql.DB, joke string) {
	var stmt *sql.Stmt
	var err error

	stmt, err = db.Prepare("INSERT INTO jokes(content) VALUES($1)")
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(joke)

	if err != nil {
		log.Fatal(err)
	}
}

// NerdyJokesHandler tells some nerdy jokes
func NerdyJokesHandler(w http.ResponseWriter, r *http.Request) {
	echoReq := r.Context().Value("echoRequest").(*alexa.EchoRequest)

	switch echoReq.GetRequestType() {
	case "LaunchRequest":
		echoResp := launchResponse()
		handleResponse(w, echoResp)
	case "IntentRequest":
		switch echoReq.GetIntentName() {
		case "TellANerdyJoke":
			echoResp := nerdyJokeResponse(echoReq)
			handleResponse(w, echoResp)
		case "HelpReply":
			echoResp := helpReply(echoReq)
			handleResponse(w, echoResp)
		case "AMAZON.HelpIntent":
			echoResp := helpResponse(echoReq)
			handleResponse(w, echoResp)
		default:
			echoResp := unknownResponse()
			handleResponse(w, echoResp)
		}
	}
}

func unknownResponse() *alexa.EchoResponse {
	return alexa.NewEchoResponse().OutputSpeech("I'm sorry, I didn't get that. Can you say that again?").EndSession(false)
}

func launchResponse() *alexa.EchoResponse {
	return alexa.NewEchoResponse().OutputSpeech("Hi, I'm Nerdy Joker. Ask me for a nerdy joke.").EndSession(false)
}

func helpResponse(echoReq *alexa.EchoRequest) *alexa.EchoResponse {
	return alexa.NewEchoResponse().OutputSpeech("You can ask me for a nerdy joke by saying \"Tell me a joke\". Do you want a joke now?").EndSession(false)
}

// Yeses are different words for acknowledgement
var Yeses = map[string]bool{
	"yes":  true,
	"sure": true,
}

func helpReply(echoReq *alexa.EchoRequest) *alexa.EchoResponse {
	want, err := echoReq.GetSlotValue("Want")

	if err != nil {
		return unknownResponse()
	}

	_, ok := Yeses[strings.ToLower(want)]

	if ok {
		return nerdyJokeResponse(echoReq)
	}

	return alexa.NewEchoResponse().OutputSpeech("Alright, have a nice day!").EndSession(true)
}

func nerdyJokeResponse(echoReq *alexa.EchoRequest) *alexa.EchoResponse {
	var content string

	err := db.QueryRow("SELECT content FROM jokes WHERE jokes.id = $1", rand.Intn(len(nerdyJokes))).Scan(&content)

	if err != nil {
		log.Fatal(err)
	}

	return alexa.NewEchoResponse().OutputSpeech(content).EndSession(true)
}

func handleResponse(w http.ResponseWriter, echoResp *alexa.EchoResponse) {
	json, _ := echoResp.String()
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}
