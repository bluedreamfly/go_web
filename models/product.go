package models

import "github.com/satori/go.uuid"

type Product struct {
	Id       uuid.UUID  `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Type     int     `json:"type"`
	ImageUrl string  `json:"image_url"`
}

