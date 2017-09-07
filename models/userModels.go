package models

import(
	"time"
)

//Administrator struct
type Administrator struct {
    Email       string `json:"email" bson:"email" binding:"required"`
    Password    string `json:"password" bson:"password" binding:"required"`
    Firstname   string `json:"firstname" bson:"firstname" binding:"required"`
    Lastname    string `json:"lastname" bson:"lastname" binding:"required"`
    Address     string `json:"address" bson:"address" binding:"required"`
    Country     string `json:"country" bson:"country" binding:"required"`
    Phone       string `json:"phone" bson:"phone" binding:"required"`
    CreatedAt   time.Time `bson:"created_at"`
}

type UpdateData struct {
    Password    string `json:"password" bson:"password" binding:"required"`
    Firstname   string `json:"firstname" bson:"firstname" binding:"required"`
    Lastname    string `json:"lastname" bson:"lastname" binding:"required"`
    Address     string `json:"address" bson:"address" binding:"required"`
    Country     string `json:"country" bson:"country" binding:"required"`
    Phone       string `json:"phone" bson:"phone" binding:"required"`
}

