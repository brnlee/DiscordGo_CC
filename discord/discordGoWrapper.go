package discord

import (
	"github.com/bwmarrin/discordgo"
)

func SendMessage(discord *discordgo.Session, msg string) {
	_, err := (*discord).ChannelMessageSend(ChannelID, msg)
	if err != nil {
		println(err.Error())
	}
}
