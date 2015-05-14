package available_topics

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/topics"
	"gopkg.in/bulletind/khabar.v1/utils"
)

const falseState = "false"
const disabledState = "disabled"

type ChotaTopic map[string]string

func GetAllTopics() []string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	topics := []string{}

	db.Conn.GetCursor(
		session, db.AvailableTopicCollection, utils.M{},
	).Distinct("ident", &topics)

	return topics
}

func GetAppTopics(app_name, org string) *[]string {
	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{"app_name": app_name}
	topics := []string{}

	var topic struct {
		Ident string `bson:"ident"`
	}

	iter := db.Conn.GetCursor(
		session, db.AvailableTopicCollection, query,
	).Select(utils.M{"ident": 1}).Sort("ident").Iter()

	for iter.Next(&topic) {
		topics = append(topics, topic.Ident)
	}

	return &topics
}

func GetAll(user, org string, appTopics, channels *[]string) (map[string]ChotaTopic, error) {
	topicMap := map[string]ChotaTopic{}

	for _, ident := range *appTopics {
		ct := ChotaTopic{"topic": ident}
		for _, channel := range *channels {
			ct[channel] = "true"
		}

		topicMap[ident] = ct
	}

	disabled := new(topics.Topic)

	session := db.Conn.Session.Copy()
	defer session.Close()

	query := utils.M{
		"ident": utils.M{"$in": appTopics},
		"user":  user,
		"org":   org,
	}

	pass1 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
	for pass1.Next(disabled) {
		if _, ok := topicMap[disabled.Ident]; ok {
			for _, blocked := range disabled.Channels {
				topicMap[disabled.Ident][blocked] = falseState
			}
		}
	}

	if user != db.BLANK {
		//Only execute this if the user was indeed passed.

		query["user"] = db.BLANK

		pass2 := db.Conn.GetCursor(session, db.TopicCollection, query).Iter()
		for pass2.Next(disabled) {
			if _, ok := topicMap[disabled.Ident]; ok {
				for _, blocked := range disabled.Channels {
					topicMap[disabled.Ident][blocked] = disabledState
				}
			}
		}
	}

	return topicMap, nil
}

func Get(topic string) (found *db.AvailableTopic, err error) {
	found = new(db.AvailableTopic)

	err = db.Conn.GetOne(db.AvailableTopicCollection, utils.M{"ident": topic}, found)

	if err != nil {
		return nil, err
	}

	return found, nil
}

func Insert(newTopic *db.AvailableTopic) string {
	return db.Conn.Insert(db.AvailableTopicCollection, newTopic)
}

func Delete(doc *utils.M) error {
	return db.Conn.Delete(db.AvailableTopicCollection, *doc)
}
