package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elvisgastelum/devsearchbot/model"
)

// Sections is used to fill in text from Google API data
var Sections [3]string

// HandleMessage is function for handle the incomming messages
func HandleMessage(url, userName, text string) {
	answer := searchAnswer(text)
	response, err := dataBinding(answer).ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	post, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(response))
	if err != nil {
		log.Fatal(err)
	}
	post.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	executePost, err := client.Do(post)
	if err != nil {
		log.Fatal(err)
	}
	defer executePost.Body.Close()
}

func searchAnswer(text string) model.SearchResults {
	text = strings.Replace(text, " ", "+", -1)
	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=AIzaSyD8QNzBdjzt3ZNEbGTz4P1rSAnvDPtbrUU&cx=005033773481765961543:gti8czyzyrw&num=3&q=%s", text)
	googleClient := http.Client{
		Timeout: time.Second * 3,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Isacc Hernandez")
	res, getErr := googleClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	value := apiMessage(body)
	return value
}

func apiMessage(jsonRaw []byte) model.SearchResults {
	jsonStructure := model.SearchResults{}
	jsonErr := json.Unmarshal(jsonRaw, &jsonStructure)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return jsonStructure
}

func dataBinding(data model.SearchResults) *model.SlashCommandResponse {

	slashCommandResponse := model.SlashCommandResponse{}
	blocks := make([]map[string]interface{}, 4)

	for i := 0; i < 3; i++ {
		item := data.Items[i]

		Sections[i] = fmt.Sprintf("*<%s|%s>*\n>_%s_", item.Link, item.Title, strings.Replace(item.Snippet, "\n", " ", -1))

		blocks[i] = map[string]interface{}{
			"type": "section",
			"accessory": map[string]interface{}{
				"type": "button",
				"text": map[string]interface{}{
					"type":  "plain_text",
					"text":  "Send",
					"emoji": true,
				},
				"value": fmt.Sprintf("button_%d", i),
			},
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": Sections[i],
			},
		}

	}

	blocks[3] = map[string]interface{}{
		"type": "actions",
		"elements": []map[string]interface{}{
			{
				"type": "button",
				"text": map[string]interface{}{
					"type":  "plain_text",
					"text":  "Cancel",
					"emoji": true,
				},
				"style": "danger",
				"value": "cancel",
			},
		},
	}

	slashCommandResponse["blocks"] = blocks

	return &slashCommandResponse
}

// ButtonAction determines the response a certain button will give
func ButtonAction(action, URL string) {
	switch action {
	case "button_0":
		jsonStr := []byte(fmt.Sprintf(`{"text":"%s","response_type":"in_channel","replace_original":true,"delete_original":true}`, Sections[0]))
		post, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Fatal(err)
		}
		post.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		executePost, er := client.Do(post)
		if er != nil {
			log.Fatal(er)
		}
		defer executePost.Body.Close()
	case "button_1":
		var jsonStr = []byte(fmt.Sprintf(`{"text":"%s","response_type":"in_channel","replace_original":true,"delete_original":true}`, Sections[1]))
		post, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Fatal(err)
		}
		post.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		executePost, er := client.Do(post)
		if er != nil {
			log.Fatal(er)
		}
		defer executePost.Body.Close()
	case "button_2":
		var jsonStr = []byte(fmt.Sprintf(`{"text":"%s","response_type":"in_channel","replace_original":true,"delete_original":true}`, Sections[2]))
		post, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Fatal(err)
		}
		post.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		executePost, er := client.Do(post)
		if er != nil {
			log.Fatal(er)
		}
		defer executePost.Body.Close()
	case "cancel":
		var jsonStr = []byte(`{"text":null,"response_type":"ephemeral","replace_original":true,"delete_original":true}`)
		post, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Fatal(err)
		}
		post.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		executePost, er := client.Do(post)
		if er != nil {
			log.Fatal(er)
		}
		defer executePost.Body.Close()
	default:
		fmt.Println("entered default event")
		log.Printf("Finish case from %s in default place\n", action)
	}

}
