package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Product struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Name        string             `bson:"name" json:"name" validate:"required,min=3,max=100"`
		Description string             `bson:"description" json:"description" validate:"max=500"`
		Price       float64            `bson:"price" json:"price" validate:"required,gt=0"`
		Category    string             `bson:"category" json:"category" validate:"required,max=50"`
		Stock       int                `bson:"stock" json:"stock" validate:"gte=0"`
	}
)
