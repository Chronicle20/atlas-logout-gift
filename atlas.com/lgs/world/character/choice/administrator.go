package choice

import (
	"gorm.io/gorm"
	"time"
)

func createChoice(db *gorm.DB) func(characterId uint32, option1 uint32, option2 uint32, option3 uint32) (Model, error) {
	return func(characterId uint32, option1 uint32, option2 uint32, option3 uint32) (Model, error) {
		c := &entity{
			CharacterId: characterId,
			Option1:     option1,
			Option2:     option2,
			Option3:     option3,
			CreatedAt:   time.Now(),
		}
		err := db.Create(c).Error
		if err != nil {
			return Model{}, err
		}
		return makeChoice(*c)
	}
}

func updateChoice(db *gorm.DB) func(characterId uint32) error {
	return func(characterId uint32) error {
		cs, err := getByCharacterId(characterId)(db)()
		if err != nil {
			return err
		}
		c := cs[0]
		c.ChosenAt = time.Now()
		return db.Save(&c).Error
	}
}
