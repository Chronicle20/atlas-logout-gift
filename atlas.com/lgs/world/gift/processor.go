package gift

import (
	"atlas-lgs/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetAll(l logrus.FieldLogger, db *gorm.DB) []Model {
	ms, err := database.ModelSliceProvider[Model, entity](db)(getAll(), makeGift)()
	if err != nil {
		l.WithError(err).Errorf("There was an error retrieving gifts")
		return make([]Model, 0)
	}
	return ms
}

func GetForWorld(l logrus.FieldLogger, db *gorm.DB) func(worldId int32) ([]Model, error) {
	return func(worldId int32) ([]Model, error) {
		ms, err := database.ModelSliceProvider[Model, entity](db)(getByWorldId(worldId), makeGift)()
		if err != nil {
			l.WithError(err).Errorf("There was an error retrieving gifts for world [%d]", worldId)
			return make([]Model, 0), err
		}
		return ms, nil
	}
}
