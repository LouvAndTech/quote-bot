package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var ErrNoChannel = errors.New("Your server does not have a citation channel")
var ErrNoCitation = errors.New("No Citation found")

/*====== Citation ======*/
//Format the citation into a string
func formatCitation(citation string, author string, date string) string {
	var t time.Time
	if date == "NONE" {
		t = time.Now().In(cstParis)
	} else {
		t, _ = time.Parse("2-1-2006", date)
	}
	strings.ReplaceAll(citation, "*", "_")
	output := fmt.Sprintf(">>> *%s*\n**%s** | *%s*", citation, author, t.Format("02.01.2006"))
	log.Printf("Format citation : %s", output)
	return output
}

// Use a button to citationised a message
func citationisation(i *discordgo.InteractionCreate) string {
	//Convert the data to something usable
	data := i.ApplicationCommandData()
	//We need to get the key to wich message we want to use
	keys := make([]string, len(data.Resolved.Messages))
	iteration := 0
	for k := range data.Resolved.Messages {
		keys[iteration] = k
		iteration++
	}
	// We get every data we need from the API
	cit := data.Resolved.Messages[keys[0]].Content
	user := fmt.Sprintf("<@!%s>", data.Resolved.Messages[keys[0]].Author.ID)
	date := data.Resolved.Messages[keys[0]].Timestamp.Format("2-1-2006")
	// Format the output
	output := formatCitation(cit, user, date)
	return output
}

func findCitationChannelID(session *discordgo.Session, guildID string) (string, error) {
	// Get all channels for the guild
	channels, err := session.GuildChannels(guildID)
	if err != nil {
		return "", err
	}

	// Loop through all channels
	for _, channel := range channels {
		// Check if the channel name contains "citation" (case-insensitive)
		if strings.Contains(strings.ToLower(channel.Name), "citation") {
			return channel.ID, nil
		}
	}

	// Return empty string if no matching channel is found
	return "", ErrNoChannel
}

func getBotMessagesInChannel(s *discordgo.Session, channelID string) ([]*discordgo.Message, error) {
	var botMessages []*discordgo.Message
	// Get the bot's ID from the session state
	botID := s.State.User.ID

	// Discord allows fetching up to 100 messages at a time, so we'll need to paginate
	limit := 100
	var lastMessageID string

	for {
		// Fetch messages with an optional "before" message ID for pagination
		messages, err := s.ChannelMessages(channelID, limit, lastMessageID, "", "")
		if err != nil {
			return nil, err
		}

		// If there are no more messages, we are done
		if len(messages) == 0 {
			break
		}

		// Filter messages sent by the bot
		for _, message := range messages {
			if message.Author != nil && message.Author.ID == botID {
				botMessages = append(botMessages, message)
			}
		}

		// Set the ID of the last message to paginate
		lastMessageID = messages[len(messages)-1].ID
	}

	if len(botMessages) == 0 {
		return botMessages, ErrNoCitation
	}

	return botMessages, nil
}

func getCitation(session *discordgo.Session, channel_id string, guildID string) (string, error) {
	output := ""
	messageList, err := getBotMessagesInChannel(s, channel_id)

	if err != nil {
		return "", err
	}

	// Pick a random element
	randomIndex := rand.Intn(len(messageList)) // random index between 0 and len(list)-1
	output = messageList[randomIndex].Content
	output += fmt.Sprintf("\nLink: https://discord.com/channels/%s/%s/%s", guildID, messageList[randomIndex].ChannelID, messageList[randomIndex].ID)
	return output, nil
}
