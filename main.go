package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func onMessage(s *discordgo.Session, m *discordgo.Message) {

}

func main() {
	discord, err := discordgo.New("Bot " + "NzA5MzMxMzM5MTkwNzMwNzgz.XrkWSg.UkG_W99GydwlcTki_aGMH4q2DLY")
	if err != nil {
		panic(err)
	}
	discord.AddHandler(onMessage)

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
