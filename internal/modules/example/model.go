package example

import "time"

type Example struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Status      string `gorm:"not null;default:active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Example) TableName() string {
	return "examples"
}
