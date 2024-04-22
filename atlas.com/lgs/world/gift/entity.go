package gift

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&entity{})
}

type entity struct {
	ID           uint32 `gorm:"primaryKey;autoIncrement;not null"`
	WorldId      int32  `gorm:"not null;default=-1"`
	ItemId       uint32 `gorm:"not null;default=0"`
	SerialNumber uint32 `gorm:"not null;default=0"`
	Weight       uint32 `gorm:"not null;default=0"`
}

func (e entity) TableName() string {
	return "gifts"
}
