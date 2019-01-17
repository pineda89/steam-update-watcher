package main

import (
	"crypto/tls"
	"github.com/thoj/go-ircevent"
	"gopkg.in/telegram-bot-api.v4"
	"os"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const (
	ircserver = "irc.freenode.net:7000"
	ircchannel = "#steamdb-announce"

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
	ircnick1 := getRandomName()
	irccon := irc.IRC(ircnick1, "IRCTestSSL")
	irccon.VerboseCallbackHandler = false
	irccon.Debug = false
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(ircchannel) })
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		manageIrcMessage(e.Message())
	})
	err := irccon.Connect(ircserver)
	if err != nil {
		log.Fatalf("Err %s", err )
		return
	}
	irccon.Loop()
}

func manageIrcMessage(message string) {
	log.Infoln("message in irc", message)

	// Here must be your filter, you can edit the source code, or use a regexp like:
	// .*(App)+.*(12578080)+.*(PLAYERUNKNOWN)+.*(BATTLEGROUNDS)

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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func getRandomName() string {
	b := make([]rune, 8)
	for i := range b {
		randata := rand.Intn(int(time.Now().UnixNano()%100000))%(len(letters)-1)
		if randata == 0 {
			randata = 1
		}
		b[i] = letters[randata]
	}
	return string(b)
}
