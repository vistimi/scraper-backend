package mongodb

import (
	"dressme-scrapper/src/utils"

	"gopkg.in/mgo.v2"
)

func Connect() *mgo.Session {

	uri := utils.DotEnvVariable("MONGODB_URI")
	session, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}
	return session
}
