package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/core"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/locks"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v3"
)

type Locks struct {
	gottp.BaseHandler
}

func (self *Locks) Post(request *gottp.Request) {
	channelIdent := request.GetArgument("channel").(string)
	topicIdent := request.GetArgument("ident").(string)

	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	lock := new(db.Locks)
	lock.PrepareSave()

	lock.Channels = []string{channelIdent}

	request.ConvertArguments(lock)
	lock.Topic = topicIdent

	if !utils.ValidateAndRaiseError(request, lock) {
		return
	}

	if locks.IsLocked(lock.Organization, lock.Topic, channelIdent, lock.Enabled) {
		request.Raise(gottp.HttpError{http.StatusConflict,
			"Already Exists."})
		return
	}

	if locks.IsLocked(lock.Organization, lock.Topic, channelIdent, !lock.Enabled) {
		request.Raise(gottp.HttpError{http.StatusConflict,
			"Already Set to Opposite. Please delete it and retry."})
		return
	}

	err, existingObj := locks.Get(lock.Organization, lock.Topic, lock.Enabled)

	if err != nil {
		if err == mgo.ErrNotFound {
			locks.Insert(lock)
			request.Write(utils.R{StatusCode: http.StatusCreated, Data: nil,
				Message: "Created"})
			return
		}
		request.Raise(gottp.HttpError{http.StatusInternalServerError,
			"Unable to complete db operation."})
		return
	}

	err = locks.AddChannel(existingObj.Topic, channelIdent, existingObj.Organization, existingObj.Enabled)

	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{http.StatusInternalServerError,
			"Unable to complete db operation."})
		return
	}

	request.Write(utils.R{
		Data:       nil,
		Message:    "true",
		StatusCode: http.StatusNoContent,
	})

	return

}

func (self *Locks) Delete(request *gottp.Request) {
	channelIdent := request.GetArgument("channel").(string)
	topicIdent := request.GetArgument("ident").(string)

	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	lock := new(db.Locks)
	request.ConvertArguments(lock)
	lock.Topic = topicIdent

	if !locks.IsLocked(lock.Organization, lock.Topic, channelIdent, lock.Enabled) {
		request.Raise(gottp.HttpError{http.StatusNotFound,
			"Does not Exists."})
		return
	}

	err := locks.RemoveChannel(lock.Topic, channelIdent, lock.Organization, lock.Enabled)

	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError,
			"Unable to complete db operation."})
		return
	}

	request.Write(utils.R{StatusCode: http.StatusNoContent, Data: nil,
		Message: "NoContent."})
	return

}
