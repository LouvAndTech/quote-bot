package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "MTA5MDU2NTE1MzQ5NTQ1MzcxNg.GjZHmG.YPX9EJVCVVdZgedoNFxm5FfM4a1_BUgGt__VHc", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

/* === Global variables === */
var s *discordgo.Session

var cstParis, _ = time.LoadLocation("Europe/Paris")

// init is called before main
func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	// Retrieve the citation channel ID from the environment variable
	//var citationChannelID := os.Getenv("CITATION_CHANNEL_ID")

}

/* === slash command initialisation === */
var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "citationisation",
			Description: "Auto format citation",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "citation",
					Description: "The text of the citation",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "author",
					Description: "The author of the citation",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date",
					Description: "if empty today date, else enter with format *day-month-year* ",
					Required:    false,
				},
			},
		},
		{
			Name:        "citation",
			Description: "Give you a random citation",
		},
		{
			Name: "instant_quote",
			Type: 3,
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"citationisation": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			citation := i.ApplicationCommandData().Options[0].StringValue()
			auteur := i.ApplicationCommandData().Options[1].StringValue()
			date := "NONE"
			if len(i.ApplicationCommandData().Options) >= 3 {
				date = i.ApplicationCommandData().Options[2].StringValue()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: formatCitation(citation, auteur, date),
				},
			})
		},
		"citation": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: sendCitation(),
				},
			})
		},
		"instant_quote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msg := citationisation(i)
			s.ChannelMessageSend("978337864301576213", msg)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
