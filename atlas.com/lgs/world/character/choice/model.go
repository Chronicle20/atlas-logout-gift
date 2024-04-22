package choice

import "time"

type Model struct {
	id          uint32
	characterId uint32
	option1     uint32
	option2     uint32
	option3     uint32
	createdAt   time.Time
	chosenAt    time.Time
}
