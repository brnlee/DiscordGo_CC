package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func SendMessage(discord *discordgo.Session, msg string) {
	_, err := (*discord).ChannelMessageSend(ChannelID, msg)
	if err != nil {
		println(err.Error())
	}
}

func SendComplexMessage(discord *discordgo.Session, content string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s cannot be opened\n> ", filename)
		return
	}
	defer file.Close()

	extension := filepath.Ext(filename)
	contentType := mime.TypeByExtension(extension)
	newFileName := uuid.New().String() + extension

	goFile := discordgo.File{Name: newFileName, Reader: file, ContentType: contentType}
	msg := discordgo.MessageSend{Content: content, File: &goFile}

	_, err = (*discord).ChannelMessageSendComplex(ChannelID, &msg)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("Uploaded %s as %s", filename, newFileName)
}

func DownloadFile(filename string, url string) {
	outfile, err := os.Create(filename)
	if err != nil {
		println(err.Error())
		return
	}
	defer outfile.Close()
	println(outfile.Name())

	resp, err := http.Get(url)
	if err != nil {
		println(err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(outfile, resp.Body)
	if err != nil {
		println(err.Error())
		return
	}
}
