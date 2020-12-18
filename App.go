package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"gopkg.in/telegram-bot-api.v4"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	ENVVAR_TELEGRAMTOKEN = "TELEGRAMTOKEN"
	ENVVAR_TELEGRAMCHATID = "TELEGRAMCHATID"
	ENVVAR_REGEXP = "REGEXP"
)

var (
	telegramToken string
	telegramChatId int64
	telegrambot *tgbotapi.BotAPI
	compiledRegexp *regexp.Regexp
)

func main() {
	loadVars()
	initTelegram()
	connectIrcAndWaitForMessages()
}

func connectIrcAndWaitForMessages() {
	config,_ := websocket.NewConfig("wss://steamdb.info/api/realtime/", "http://localhost.localdomain/")
	ws, err := websocket.DialConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		msg := make([]byte, 512)
		_, err = ws.Read(msg)
		if err != nil {
			log.Fatalln(err)
		}

		manageMessage(string(msg[:bytes.Index(msg, []byte{0, 0, 0, 0})]))
	}
}

func manageMessage(message string) {
	log.Infoln("message received", message)

	// Here must be your filter, you can edit the source code, or use a regexp like:
	// .*(12578080)+.*PLAYERUNKNOWN*

	if compiledRegexp.MatchString(message) {
		log.Infoln("message matchs the regexp")
		telegrambot.Send(tgbotapi.NewMessage(telegramChatId, message))
		log.Infoln("message sended to telegram", message)
	}
}

func initTelegram() {
	var err error
	telegrambot, err = tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalln("error creating bot telegram", err)
		panic(err)
	}

	log.Infoln("telegram bot created")
}

func loadVars() {
	rand.Seed(time.Now().UTC().UnixNano())
	var err error

	telegramToken = getEnvVar(ENVVAR_TELEGRAMTOKEN)
	log.Infoln("telegramToken", telegramToken)
	telegramChatId, err = strconv.ParseInt(getEnvVar(ENVVAR_TELEGRAMCHATID), 10, 64)
	if err != nil {
		log.Fatalln("TELEGRAMCHATID parsing error", err)
		panic(err)
	}

	log.Infoln("telegramChatId", telegramChatId)
	compiledRegexp = regexp.MustCompile(getEnvVar(ENVVAR_REGEXP))
	log.Infoln("regexp", getEnvVar(ENVVAR_REGEXP))

	log.Infoln("env vars loaded")
}

func getEnvVar(s string) string {
	result := os.Getenv(s)
	if result == "" {
		log.Fatalln(s, "env var not informed")
		panic(s + " env var not informed")
	}

	return result
}
