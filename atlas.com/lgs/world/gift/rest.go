package gift

import "strconv"

type RestModel struct {
	Id           uint32 `json:"-"`
	WorldId      uint32 `json:"world_id"`
	ItemId       uint32 `json:"item_id"`
	SerialNumber uint32 `json:"serial_number"`
	Weight       uint32 `json:"weight"`
}

func (r RestModel) GetName() string {
	return "gifts"
}

func (r RestModel) GetID() string {
	return strconv.Itoa(int(r.Id))
}

func TransformAll(models []Model) []RestModel {
	rms := make([]RestModel, 0)
	for _, m := range models {
		rms = append(rms, Transform(m))
	}
	return rms
}

func Transform(model Model) RestModel {
	rm := RestModel{
		Id:           model.id,
		ItemId:       model.itemId,
		SerialNumber: model.serialNumber,
		Weight:       model.weight,
	}
	return rm
}
