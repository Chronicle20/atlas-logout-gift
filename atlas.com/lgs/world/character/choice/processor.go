package choice

import (
	"atlas-lgs/database"
	"atlas-lgs/model"
	"atlas-lgs/world/gift"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"sort"
	"time"
)

func GetForCharacter(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		ms, err := model.First(database.ModelSliceProvider[Model, entity](db)(getByCharacterId(characterId), makeChoice))
		if err != nil {
			if errors.Is(err, model.EmptyResult) {
				return Model{}, err
			}
			l.WithError(err).Errorf("There was an error retrieving choices for character [%d]", characterId)
			return Model{}, err
		}
		return ms, nil
	}
}

func GetForWorldCharacter(l logrus.FieldLogger, db *gorm.DB) func(worldId uint32, characterId uint32) (Model, error) {
	return func(worldId uint32, characterId uint32) (Model, error) {
		c, err := GetForCharacter(l, db)(characterId)
		if err != nil {
			if errors.Is(err, model.EmptyResult) {
				//user gifts don't exist
				c = ProduceChoices(l, db)(worldId, characterId)
				return c, nil
			}
			return Model{}, err
		}

		age := time.Since(c.createdAt)
		threeDays := 3 * 24 * time.Hour

		if age > threeDays {
			//	user has not had a gift within configurable days
			//		generate a set of gifts from pool
			//	gifts were generated more than # of days
			//		generate new gift set
			c = ProduceChoices(l, db)(worldId, characterId)
			return c, nil
		}

		if c.chosenAt.IsZero() {
			//	gifts are still fresh
			//		great, return them
			return c, nil
		}
		return Model{}, nil
	}
}

type ByWeight []gift.Model

func (a ByWeight) Len() int { return len(a) }

func (a ByWeight) Less(i, j int) bool { return a[i].Weight() > a[j].Weight() }

func (a ByWeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func ProduceChoices(l logrus.FieldLogger, db *gorm.DB) func(worldId uint32, characterId uint32) Model {
	return func(worldId uint32, characterId uint32) Model {
		gs, err := gift.GetForWorld(l, db)(int32(worldId))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			gs, err = gift.GetForWorld(l, db)(-1)
			if err != nil {
				l.WithError(err).Errorf("Unable to get gift choices.")
				return Model{}
			}
		}

		if err != nil {
			l.WithError(err).Errorf("Unable to get gift for world [%d].", worldId)
			return Model{}
		}

		// Sort options by weight
		sort.Slice(gs, func(i, j int) bool {
			return gs[i].Weight() > gs[j].Weight()
		})

		// Create a slice to store selected options
		selected := make([]gift.Model, 0, 3)

		// Create a map to keep track of selected options to ensure uniqueness
		selectedMap := make(map[uint32]bool)

		// Select unique options based on their weights
		for len(selected) <= 3 && len(gs) > 0 {
			// Calculate total weight
			totalWeight := int32(0)
			for _, opt := range gs {
				totalWeight += int32(opt.Weight())
			}

			// Generate a random number in the range [0, totalWeight)
			r := rand.Int31n(totalWeight)

			// Select an option based on the random number
			cumulativeWeight := int32(0)
			var opt gift.Model
			for _, opt = range gs {
				cumulativeWeight += int32(opt.Weight())
				if r <= cumulativeWeight {
					if !selectedMap[opt.SerialNumber()] {
						selected = append(selected, opt)
						selectedMap[opt.SerialNumber()] = true
						break
					}
				}
			}

			// Remove the selected option from the options slice to prevent duplicates
			ngs := make([]gift.Model, 0, len(gs)-1)
			for _, g := range gs {
				if g.SerialNumber() != opt.SerialNumber() {
					ngs = append(ngs, g)
				}
			}
			gs = ngs
		}

		c, err := createChoice(db)(characterId, selected[0].SerialNumber(), selected[1].SerialNumber(), selected[2].SerialNumber())
		if err != nil {
			l.WithError(err).Errorf("There was an error creating choices for world [%d] character [%d]", worldId, characterId)
			return Model{}
		}
		return c
	}
}

func ProcessChoice(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) error {
	return func(characterId uint32) error {
		return updateChoice(db)(characterId)
	}
}
