package pe_data_models

import "go.mongodb.org/mongo-driver/bson/primitive"

type KthkMsgs struct {
	Mid      primitive.ObjectID `bson:"_id" json:"_id"`
	UserId   string             `bson:"USER_ID" json:"USER_ID"`
	RoomName string             `bson:"ROOM_NAME" json:"ROOM_NAME"`
	Channel  string             `bson:"CHANNEL" json:"CHANNEL"`
	Content  string             `bson:"CONTENT" json:"CONTENT"`
	Ts       uint64             `bson:"TS" json:"TS"`
	Iid      uint64             `bson:"IID" json:"IID"`
	Rrating  uint32             `bson:"R_RATING" json:"R_RATING"`
}
