package discord

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SendMessage(discord *discordgo.Session, msg string) {
	if strings.HasSuffix(msg, "\n") == false {
		msg += "\n"
	}
	_, err := (*discord).ChannelMessageSend(ChannelID, encrypt(msg))
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
	msg := discordgo.MessageSend{Content: encrypt(content), File: &goFile}

	_, err = (*discord).ChannelMessageSendComplex(ChannelID, &msg)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("Uploaded %s as %s\n", filename, newFileName)
}

func DownloadFile(filename string, url string) {
	outfile, err := os.Create(filename)
	if err != nil {
		println(err.Error())
		return
	}
	defer outfile.Close()

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
	fmt.Printf("%s has been saved to this progrgam's working directory\n", filename)
}

// Encrypt and Decrypt functions copied from https://gist.github.com/jyap808/8250124
func encrypt(msg string) string {
	encryptionPassphrase := []byte(Secret)
	encryptionText := msg
	encryptionType := "PGP SIGNATURE"

	encbuf := bytes.NewBuffer(nil)
	w, err := armor.Encode(encbuf, encryptionType, nil)
	if err != nil {
		log.Fatal(err)
	}

	plaintext, err := openpgp.SymmetricallyEncrypt(w, encryptionPassphrase, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	message := []byte(encryptionText)
	_, err = plaintext.Write(message)

	plaintext.Close()
	w.Close()
	//fmt.Printf("Encrypted:\n%s\n", encbuf)
	return encbuf.String()
}

func Decrypt(msg string) string {
	decbuf := bytes.NewBuffer([]byte(msg))
	result, err := armor.Decode(decbuf)
	if err != nil {
		log.Fatal(err)
	}

	md, err := openpgp.ReadMessage(result.Body, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return []byte(Secret), nil
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	//fmt.Printf("Decrypted:\n%s\n", string(bytes))
	return string(bytes)
}
