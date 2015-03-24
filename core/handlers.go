package core

import "github.com/changer/khabar/dbapi/pending"

const (
	EMAIL = "email"
	SMS   = "sms"
)

var ChannelMap = map[string]func(*pending.PendingItem, string, map[string]interface{}){
	EMAIL: emailHandler,
}
