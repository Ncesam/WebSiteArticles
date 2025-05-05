package mongo


type Config struct {
	Id           int64  `bson:"id"`
	UserId       int64  `bson:"userId"`
	Name         string `bson:"name"`
	Prompt       string `bson:"prompt"`
	RefreshTokenDTF string `bson:"refreshTokenDTF"`
	RefreshTokenVC string `bson:"refreshTokenVC"`
	Email        string `bson:"email"`
	Password     string `bson:"password"`
	Delay        int64  `bson:"delay"`
}

