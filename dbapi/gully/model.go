package gully

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/dbapi"
)

const (
	GullyCollection = "gullys"
)

type Gully struct {
	db.BaseModel   `bson:",inline"`
	UserID         string                 `json:"user_id" bson:"user_id"`
	OrganizationID string                 `json:"org_id" bson:"org_id"`
	ApplicationID  string                 `json:"app_id" bson:"app_id"`
	GullyData      map[string]interface{} `json:"channel_data" bson:"channel_data" required:"true"`
	Ident          string                 `json:"ident" bson:"ident" required:"true"`
}

func (self *Gully) IsValid(op_type int) bool {
	if (len(self.UserID) == 0) && (len(self.OrganizationID) == 0) && (len(self.ApplicationID) == 0) {
		return false
	}

	if len(self.Ident) == 0 {
		return false
	}

	if op_type == dbapi.INSERT_OPERATION {
		if len(self.GullyData) == 0 {
			return false
		}

	}

	return true
}
