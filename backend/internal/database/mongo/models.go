package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type ConfigSchema struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"` // MongoDB использует ObjectID для идентификаторов
	Prompt  string             `bson:"prompt"`
	UserID  int64              `bson:"userId"`
	Account Account            `bson:"account"`
}

type Account struct {
	ID       int64  `bson:"id"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
