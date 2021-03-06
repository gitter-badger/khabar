package user_locale

import (
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/utils"
)

func Get(user string) (userLocale *db.UserLocale, err error) {
	userLocale = new(db.UserLocale)
	err = db.Conn.GetOne(db.UserLocaleCollection, utils.M{"user": user},
		userLocale)
	if err != nil {
		return nil, err
	}
	return
}

func Insert(userLocale *db.UserLocale) string {
	return db.Conn.Insert(db.UserLocaleCollection, userLocale)
}

func Update(user string, doc *utils.M) error {
	return db.Conn.Update(db.UserLocaleCollection, utils.M{"user": user},
		utils.M{
			"$set": *doc,
		})
}
