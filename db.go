package main

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SubmittedData struct {
	Type        string `bson:"_type"` // data type
	CreatedDate time.Time
	UpdateDate  time.Time

	// Optional fields
	UserID          string
	ModerationFlags []string
	data            *bson.M
}

func dbConnect(url string) *mgo.Session {
	log.Printf("Connecting to mongo @ %s", url)
	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	return session
}

func findSubmittedData(c *mgo.Collection, ty string, skip int, limit int) ([]SubmittedData, error) {
	var a []SubmittedData
	err := c.Find(bson.M{"_type": ty}).Sort("-UpdateDate").Skip(skip).Limit(limit).All(&a)
	return a, err
}
