package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID         primitive.ObjectID `bson:"_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Table_id   *string            `json:"table_id" `
	Order_id   string             `json:"order_id" validate:"required"`
}
