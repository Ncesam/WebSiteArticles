package mongo

import "go.mongodb.org/mongo-driver/v2/bson"



type Config struct {
	Id           bson.ObjectID `bson:"_id,omitempty"`
	UserId       int64  `bson:"userId"`
	Name         string `bson:"name"`
	Prompt       string `bson:"prompt"`
	RefreshTokenDTF string `bson:"refreshTokenDTF"`
	RefreshTokenVC string `bson:"refreshTokenVC"`
	Email        string `bson:"email"`
	Password     string `bson:"password"`
	Delay        int64  `bson:"delay"`
}

