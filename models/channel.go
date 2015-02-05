package models

import (
	"github.com/parthdesai/sc-notifications/db"
)

type Channel struct {
	db.BaseModel   `bson:",inline"`
	UserID         string                 `json:"user_id" bson:"user_id"`
	OrganizationID string                 `json:"org_id" bson:"org_id"`
	ApplicationID  string                 `json:"app_id" bson:"app_id"`
	ChannelData    map[string]interface{} `json:"channel_data" bson:"channel_data" required:"true"`
	Ident          string                 `json:"ident" bson:"ident" required:"true"`
}

func (self *Channel) IsValid() bool {
	if (len(self.UserID) == 0) && (len(self.OrganizationID) == 0) && (len(self.ApplicationID) == 0) {
		return false
	}

	if len(self.Ident) == 0 {
		return false
	}

	if len(self.ChannelData) == 0 {
		return false
	}

	return true
}

func (self *Channel) GetFromDatabase(dbConn *db.MConn) bool {
	return dbConn.Get("channels", db.M{"app_id": self.ApplicationID,
		"org_id": self.OrganizationID, "user_id": self.UserID, "ident": self.Ident}).Next(self)
}

func (self *Channel) DeleteFromDatabase(dbConn *db.MConn) error {
	return dbConn.Delete("channels", db.M{"app_id": self.ApplicationID,
		"org_id": self.OrganizationID, "user_id": self.UserID, "ident": self.Ident})
}
