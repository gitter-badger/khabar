package core

import (
	"bytes"

	"log"
	"net/http"

	"github.com/bulletind/khabar/dbapi/pending"
	"gopkg.in/simversity/gottp.v2/utils"
)

const PARSE_URL = "https://api.parse.com/1/push"

func parseHandler(
	item *pending.PendingItem,
	text string,
	settings map[string]interface{},
) {
	log.Println("Sending parse...")

	application_id, ok := settings["parse_application_id"].(string)
	if !ok {
		log.Println("parse_application_id is a required parameter.")
		return
	}

	api_key, ok := settings["parse_rest_api_key"].(string)
	if !ok {
		log.Println("parse_rest_api_key is a required parameter.")
		return
	}

	body := map[string]interface{}{}
	body["alert"] = text
	body["entity"] = item.Entity
	body["organization"] = item.Organization
	body["app_name"] = item.AppName
	body["topic"] = item.Topic
	body["created_on"] = item.CreatedOn

	var jsonStr = utils.Encoder(&body)

	req, err := http.NewRequest("POST", PARSE_URL, bytes.NewBuffer(jsonStr))

	req.Header.Set("X-Parse-Application-Id", application_id)
	req.Header.Set("X-Parse-REST-API-Key", api_key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Parse push Status:", resp.Status)
}
