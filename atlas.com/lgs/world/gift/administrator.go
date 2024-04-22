package gift

import "gorm.io/gorm"

func BulkCreateGifts(db *gorm.DB, gifts []Model) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	for _, gm := range gifts {
		g := &entity{
			WorldId:      gm.worldId,
			ItemId:       gm.itemId,
			SerialNumber: gm.serialNumber,
			Weight:       gm.weight,
		}

		err := tx.Create(g).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
