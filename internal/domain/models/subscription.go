package models

import "github.com/google/uuid"

type Subscription struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name" validate:"required,min=2"`
	Price       int       `json:"price" validate:"required,gt=0"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	StartDate   string    `json:"start_date" validate:"required,mm_yyyy"`
	EndDate     *string   `json:"end_date,omitempty" validate:"omitempty,mm_yyyy"`
}

type ListSubsParams struct {
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	ServiceName string    `json:"service_name" validate:"required,min=2"`
	StartDate   string    `json:"start_date" validate:"omitempty,mm_yyyy"`
	EndDate     string    `json:"end_date" validate:"omitempty,mm_yyyy"`
}
