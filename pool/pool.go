package pool

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var session *mgo.Session

func NewConnection() {
	var err error
	session, err = mgo.Dial("localhost")
	if err != nil {
		log.Fatal("[pool.go] Init - cannot dial mongo : ", err)
	}

	log.Println("[pool] NewConnection Success")
}

func Close() {
	session.Close()
}

func GetSession() *mgo.Session {
	return session.Copy()
}
