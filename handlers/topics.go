package handlers

import (
	"log"
	"net/http"

	"github.com/bulletind/khabar/core"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/available_topics"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v2"
)

type TopicChannel struct {
	gottp.BaseHandler
}

func (self *TopicChannel) Delete(request *gottp.Request) {
	intopic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)

	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	intopic.Ident = request.GetArgument("ident").(string)

	request.ConvertArguments(intopic)

	topic, err := topics.Get(
		intopic.User,
		intopic.Organization,
		intopic.Ident,
	)

	if err != nil && err != mgo.ErrNotFound {
		log.Println(err)
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to fetch data, Please try again later.",
		})

		return

	}

	var hasData bool

	if topic == nil {
		log.Println("Creating new document")
		intopic.AddChannel(channelIdent)

		intopic.PrepareSave()
		if !intopic.IsValid(db.INSERT_OPERATION) {
			request.Raise(gottp.HttpError{
				http.StatusBadRequest,
				"Atleast one of the user or org must be present.",
			})

			return
		}

		if !utils.ValidateAndRaiseError(request, intopic) {
			log.Println("Validation Failed")
			return
		}

		topic = intopic

	} else {
		hasData = true

		for _, ident := range topic.Channels {
			if ident == channelIdent {
				request.Raise(gottp.HttpError{
					http.StatusConflict,
					"You have already unsubscribed this channel",
				})

				return
				break
			}
		}

		topic.AddChannel(channelIdent)
	}

	if hasData {
		err = topics.Update(
			topic.User,
			topic.Organization,
			topic.Ident,
			&utils.M{"channels": topic.Channels},
		)

		if err != nil {
			log.Println("Error while inserting document :" + err.Error())
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Internal server error.",
			})

			return
		} else {
			request.Write(utils.R{
				Data:       nil,
				Message:    "true",
				StatusCode: http.StatusNoContent,
			})

			return
		}
	} else {
		log.Println("Successfull call: Inserting document")
		topics.Insert(topic)
		request.Write(utils.R{
			Data:       nil,
			Message:    "true",
			StatusCode: http.StatusNoContent,
		})

		return
	}
}

func (self *TopicChannel) Post(request *gottp.Request) {
	topic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)
	topic.Ident = request.GetArgument("ident").(string)

	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	request.ConvertArguments(topic)

	topic, err := topics.Get(
		topic.User,
		topic.Organization,
		topic.Ident,
	)

	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})

		} else {
			request.Write(utils.R{
				Data:       nil,
				Message:    "true",
				StatusCode: http.StatusNoContent,
			})
		}

		return
	}

	topic.RemoveChannel(channelIdent)
	log.Println(topic.Channels)

	if len(topic.Channels) == 0 {
		log.Println("Deleting from database, since channels are now empty.")
		err = topics.Delete(

			&utils.M{
				"org":   topic.Organization,
				"user":  topic.User,
				"ident": topic.Ident,
			},
		)

	} else {
		log.Println("Updating...")

		err = topics.Update(
			topic.User, topic.Organization,
			topic.Ident, &utils.M{"channels": topic.Channels},
		)
	}

	if err != nil {
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to delete.",
		})

		return
	}

	request.Write(utils.R{
		Data:       nil,
		Message:    "true",
		StatusCode: http.StatusNoContent,
	})

	return
}

type Topics struct {
	gottp.BaseHandler
}

func (self *Topics) Get(request *gottp.Request) {
	var args struct {
		Organization string `json:"org" required:"true"`
		AppName      string `json:"app_name" required:"true"`
		User         string `json:"user"`
	}

	request.ConvertArguments(&args)

	if !utils.ValidateAndRaiseError(request, args) {
		log.Println("Validation Failed")
		return
	}

	channels := []string{}
	for ident, _ := range core.ChannelMap {
		channels = append(channels, ident)
	}

	iter, err := available_topics.GetAll(args.User, args.AppName, args.Organization, channels)

	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})

		} else {
			request.Raise(gottp.HttpError{
				http.StatusNotFound,
				"Not Found.",
			})
		}

		return
	}

	ret := []available_topics.ChotaTopic{}
	for _, singleRet := range iter {
		ret = append(ret, singleRet)
	}

	request.Write(ret)
	return
}

func (self *Topics) Post(request *gottp.Request) {
	newTopic := new(db.AvailableTopic)

	request.ConvertArguments(newTopic)

	newTopic.PrepareSave()

	if !utils.ValidateAndRaiseError(request, newTopic) {
		log.Println("Validation Failed")
		return
	}

	if _, err := available_topics.Get(newTopic.Ident); err == nil {
		request.Raise(gottp.HttpError{
			http.StatusConflict,
			"Topic already exists"})
		return
	} else {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})
			return
		}
	}

	available_topics.Insert(newTopic)

	request.Write(utils.R{
		StatusCode: http.StatusCreated,
		Data:       newTopic.Id,
		Message:    "Created",
	})
	return
}

type Topic struct {
	gottp.BaseHandler
}

func (self *Topics) Delete(request *gottp.Request) {
	var args struct {
		Ident string `json:"ident" required:"true"`
	}

	request.ConvertArguments(&args)

	if !utils.ValidateAndRaiseError(request, args) {
		log.Println("Validation Failed")
		return
	}

	err := topics.DeleteTopic(args.Ident)
	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to delete.",
		})
		return
	}

	request.Write(utils.R{
		Data:       nil,
		Message:    "true",
		StatusCode: http.StatusNoContent,
	})
	return
}
