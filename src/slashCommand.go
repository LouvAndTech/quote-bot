package main

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

/* === slash command initialisation === */
var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "quotization",
			Description: "Auto format quote",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "quote",
					Description: "The text of the quote",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "author",
					Description: "The author of the quote",
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
		/*{
			Name:        "correct_citationisation",
			Description: "Correct a citation",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID of the message to correct",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "citation",
					Description: "The corrected text of the citation",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "author",
					Description: "The corrected author of the citation",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date",
					Description: "Corrected date with format *day-month-year* ",
					Required:    false,
				},
			},
		},*/
		{
			Name:        "quote",
			Description: "Give you a random quote from the server",
		},
		{
			Name: "instant_quote",
			Type: 3,
		},
		{
			Name: "remove_quote",
			Type: 3,
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"quotization": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			content, flags := "An error occured", discordgo.MessageFlags(1<<6)
			// get the channel ID
			channel_id, err := findQuoteChannelID(s, i.GuildID)
			if err != nil {
				if errors.Is(err, ErrNoChannel) {
					content = err.Error()
				} else {
					log.Println(err)
				}
			}
			//Prevent the use of the commqnd in the citation channel to avoid duplicate
			if i.ChannelID == channel_id {
				citation := i.ApplicationCommandData().Options[0].StringValue()
				auteur := i.ApplicationCommandData().Options[1].StringValue()
				date := "NONE"
				if len(i.ApplicationCommandData().Options) >= 3 {
					date = i.ApplicationCommandData().Options[2].StringValue()
				}
				content = formatQuote(citation, auteur, date)
				flags = discordgo.MessageFlags(0)
			} else {
				content = "Wrong channel"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   flags,
				},
			})
		},
		/* "correct_citationisation": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Print("TEST")
		},*/
		"quote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			content, flags := "An error occured", discordgo.MessageFlags(1<<6)
			// get the channel ID
			channel_id, err := findQuoteChannelID(s, i.GuildID)
			if err != nil {
				if errors.Is(err, ErrNoChannel) {
					content = err.Error()
				} else {
					log.Println(err)
				}
			}
			//Prevent the use of the commqnd in the citation channel to avoid duplicate
			if i.ChannelID == channel_id {
				content = "You cannot use this command in the citation channel"
			} else {
				//Get the content
				content, err = getCitation(s, channel_id, i.GuildID)
				if err != nil {
					if errors.Is(err, ErrNoCitation) {
						content = err.Error()
					} else {
						log.Println(err)
					}
				} else {
					flags = discordgo.MessageFlags(0)
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   flags,
				},
			})
		},
		"instant_quote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msg := quotization(i)
			channel_id, err := findQuoteChannelID(s, i.GuildID)
			content := "An error occured"
			if err != nil {
				if errors.Is(err, ErrNoChannel) {
					content = err.Error()
				} else {
					log.Println(err)
				}
			} else {
				content = "Done"
				s.ChannelMessageSend(channel_id, msg)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})
		},
		"remove_quote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := removeQuote(s, i)
			content := "An error occured"
			if err != nil {
				if errors.Is(err, ErrPermMsg) {
					content = err.Error()
				} else {
					log.Println(err)
				}
			} else {
				content = "Message removed"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Flags:   1 << 6,
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "An unexpected error occurred. Please try again later.",
						Flags:   1 << 6,
					},
				})
			}
		}()
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
