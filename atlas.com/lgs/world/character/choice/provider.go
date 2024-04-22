package choice

import (
	"atlas-lgs/database"
	"atlas-lgs/model"
	"gorm.io/gorm"
)

func getByCharacterId(characterId uint32) database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		var results []entity
		err := db.Where(&entity{CharacterId: characterId}).Order("created_at desc").Find(&results).Error
		if err != nil {
			return model.ErrorSliceProvider[entity](err)
		}
		return model.FixedSliceProvider(results)
	}
}

func makeChoice(m entity) (Model, error) {
	r := Model{
		id:          m.ID,
		characterId: m.CharacterId,
		option1:     m.Option1,
		option2:     m.Option2,
		option3:     m.Option3,
		createdAt:   m.CreatedAt,
		chosenAt:    m.ChosenAt,
	}
	return r, nil
}
