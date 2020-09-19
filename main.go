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

var sessions map[string]chan string

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	_, exists := sessions[m.ChannelID]
	if exists { //session already exists, redirect output
		sessions[m.ChannelID] <- m.Content
	} else {
		execCommand(s, m)
	}

}

func redirectOutput(pipe *io.ReadCloser, s *discordgo.Session, m *discordgo.MessageCreate) {
	scanner := bufio.NewScanner(*pipe)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("Got from stdout  " + text)
		_, err := s.ChannelMessageSend(m.ChannelID, text)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func redirectInput(pipe *io.WriteCloser, m *discordgo.MessageCreate) {

	for s := range sessions[m.ChannelID] {
		fmt.Println("Wrote to stdin  " + s)
		_, err := (*pipe).Write([]byte(s + "\n"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func execCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	comAndArg := strings.Split(m.Content, " ")

	if comAndArg[0] == "[i]" {
		sessions[m.ChannelID] = make(chan string)
		cmd := exec.Command(comAndArg[1], comAndArg[2:]...)
		outP, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}
		inP, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}
		errP, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println(err)
		}
		go redirectInput(&inP, m)
		go redirectOutput(&outP, s, m)
		go redirectOutput(&errP, s, m)
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
		}

		//_, err = inP.Write([]byte("exit()"))
		if err != nil {
			fmt.Println(err)
		}

		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
		close(sessions[m.ChannelID])
		delete(sessions, m.ChannelID)
	} else {
		cmd := exec.Command(comAndArg[0], comAndArg[1:]...)
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			return
		}
		for i := 0; i <= (len(stdoutStderr) / 1900); i++ {
			slicedTo := (i + 1) * 1900
			if slicedTo > len(stdoutStderr) {
				slicedTo = len(stdoutStderr) - 1
			}
			_, err = s.ChannelMessageSend(m.ChannelID, string(stdoutStderr[i*1900:slicedTo]))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func main() {
	sessions = make(map[string]chan string)
	discord, err := discordgo.New("Bot " + "Token")
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
