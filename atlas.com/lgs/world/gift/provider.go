package gift

import (
	"atlas-lgs/database"
	"atlas-lgs/model"
	"gorm.io/gorm"
)

func getAll() database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		var results []entity
		err := db.Find(&results).Error
		if err != nil {
			return model.ErrorSliceProvider[entity](err)
		}
		return model.FixedSliceProvider(results)
	}
}

func getByWorldId(worldId int32) database.EntitySliceProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		return database.SliceQuery[entity](db, map[string]interface{}{"world_id": worldId})
	}
}

func makeGift(m entity) (Model, error) {
	r := NewGiftBuilder(m.ID).
		SetWorldId(m.WorldId).
		SetItemId(m.ItemId).
		SetSerialNumber(m.SerialNumber).
		SetWeight(m.Weight).
		Build()
	return r, nil
}
