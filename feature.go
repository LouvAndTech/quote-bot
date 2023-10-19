package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

/*====== Citation ======*/
//Format the citation into a string
func formatCitation(citation string, author string, date string) string {
	var t time.Time
	if date == "NONE" {
		t = time.Now().In(cstParis)
	} else {
		t, _ = time.Parse("2-1-2006", date)
	}
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

type Citation struct {
	Text   string
	Author string
}

func getCitation() Citation {
	response, err := http.Get("https://type.fit/api/quotes")
	if err != nil {
		fmt.Print(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var citations []Citation
	json.Unmarshal(responseData, &citations)
	//log.Printf("citation 1 : \ntext : %s\nauthor : %s", citations[0].Text, citations[0].Author)
	//log.Printf("nb citation : %d", len(citations))
	return citations[rand.Intn(len(citations))]
}

func sendCitation() string {
	citation := getCitation()
	output := fmt.Sprintf("> *%s*\n> **%s**", citation.Text, citation.Author)
	return output
}
