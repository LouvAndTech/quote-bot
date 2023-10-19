package main

import (
	"github.com/bwmarrin/discordgo"
)

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
