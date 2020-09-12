package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

}

func redirectOutput(pipe *io.Reader, s *discordgo.Session) {
	scanner := bufio.NewScanner(*pipe)
	for scanner.Scan() {
		scanner.Text()
	}
}

func execCommand(command string, args string, stdout chan string, stdin chan string) {
	cmd := exec.Command(command, args)
	outP, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	outP.Read()
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
