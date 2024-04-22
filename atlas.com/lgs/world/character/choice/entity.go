package choice

import (
	"gorm.io/gorm"
	"time"
)

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&entity{})
}

type entity struct {
	ID          uint32    `gorm:"primaryKey;autoIncrement;not null"`
	CharacterId uint32    `gorm:"not null;default=0"`
	Option1     uint32    `gorm:"not null;default=0"`
	Option2     uint32    `gorm:"not null;default=0"`
	Option3     uint32    `gorm:"not null;default=0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	ChosenAt    time.Time
}

func (e entity) TableName() string {
	return "choices"
}
