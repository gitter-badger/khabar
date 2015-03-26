package handlers

import (
	"log"
	"net/http"

	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v2"
)

type TopicChannel struct {
	gottp.BaseHandler
}

func (self *TopicChannel) Post(request *gottp.Request) {
	intopic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)

	//FIXME: Use some common location for this function instead of
	// handlers/gully.go.

	if !isChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	intopic.Ident = request.GetArgument("ident").(string)

	request.ConvertArguments(intopic)

	topic, err := topics.Get(
		intopic.User, intopic.AppName,
		intopic.Organization, intopic.Ident,
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
				"Atleast one of the user, org and app_name must be present.",
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
					"Channel is already a part of this Topic.",
				})

				return
				break
			}
		}

		topic.AddChannel(channelIdent)
	}

	if hasData {
		err = topics.Update(
			topic.User, topic.AppName,
			topic.Organization, topic.Ident,
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
				Message:    "NoContent",
				StatusCode: http.StatusNoContent,
			})

			return
		}
	} else {
		log.Println("Successfull call: Inserting document")
		topics.Insert(topic)
		request.Write(utils.R{
			Data:       topic.Id,
			Message:    "Created",
			StatusCode: http.StatusCreated,
		})

		return
	}
}

func (self *TopicChannel) Delete(request *gottp.Request) {
	topic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)
	topic.Ident = request.GetArgument("ident").(string)

	request.ConvertArguments(topic)

	topic, err := topics.Get(
		topic.User, topic.AppName,
		topic.Organization, topic.Ident,
	)

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

	if topic == nil {
		request.Raise(gottp.HttpError{
			http.StatusNotFound,
			"topics setting does not exists.",
		})

		return
	}

	topic.RemoveChannel(channelIdent)
	log.Println(topic.Channels)

	if len(topic.Channels) == 0 {
		log.Println("Deleting from database, since channels are now empty.")
		err = topics.Delete(

			&utils.M{
				"app_name": topic.AppName,
				"org":      topic.Organization,
				"user":     topic.User,
				"ident":    topic.Ident,
			},
		)

	} else {
		log.Println("Updating...")

		err = topics.Update(
			topic.User, topic.AppName, topic.Organization,
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
		Message:    "NoContent",
		StatusCode: http.StatusNoContent,
	})

	return
}

type Topic struct {
	gottp.BaseHandler
}

func (self *Topic) Delete(request *gottp.Request) {
	topic := new(topics.Topic)
	request.ConvertArguments(topic)
	if !topic.IsValid(db.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Atleast one of the user, org and app_name must be present.",
		})

		return
	}

	err := topics.Delete(

		&utils.M{
			"app_name": topic.AppName,
			"org":      topic.Organization,
			"user":     topic.User,
			"ident":    topic.Ident,
		},
	)

	if err != nil {
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to delete.",
		})

		return
	}

	request.Write(utils.R{Data: nil, Message: "NoContent",
		StatusCode: http.StatusNoContent})
	return
}

type Topics struct {
	gottp.BaseHandler
}

func (self *Topics) Get(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user"`
	}

	request.ConvertArguments(&args)

	all, err := topics.GetAll(args.User, args.AppName,
		args.Organization)

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

	request.Write(all)
	return
}
