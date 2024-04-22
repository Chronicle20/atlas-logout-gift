package choice

import "strconv"

type RestModel struct {
	Id      uint32 `json:"-"`
	Option1 uint32 `json:"option_1"`
	Option2 uint32 `json:"option_2"`
	Option3 uint32 `json:"option_3"`
}

func (r RestModel) GetName() string {
	return "choices"
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
		Id:      model.id,
		Option1: model.option1,
		Option2: model.option2,
		Option3: model.option3,
	}
	return rm
}
