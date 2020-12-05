
package models

import (
"go.mongodb.org/mongo-driver/bson/primitive"
"time"
)

type Stamp struct {
	TimeStamp string `json:"time_stamp" bson:"time_stamp"`
	Topic     string        `json:"topic" bson:"topic"`
}

type Podcast struct {
	Name     string  `json:"name" bson:"name"`
	PODID       string  `json:"podid" bson:"podid"`
	ImageUrl string  `json:"image_url" bson:"image_url"`
	Transcript string	`json:"transcript" bson:"transcript"`
	Topics   []Stamp `json:"topics" bson:"topics"`
}

type User struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	Sub              string             `json:"sub" bson:"sub"`
	Username         string             `bson:"username" json:"username"`
	FullName         string             `json:"full_name" bson:"full_name"`
	Email            string             `json:"email" bson:"email"`
	Password         string             `json:"password" bson:"password"`
	FavoritePodcast  []Podcast          `json:"favorite_podcast bson:"favorite_podcast"`
	InterestedTopics []string           `json:"interested_topics bson:"interested_topics"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
	Active           bool               `json:"active" bson:"active"`
	Method           string             `json:"method" bson:"method"`
	//Token			string				`json:"token" bson:"token"`
}

