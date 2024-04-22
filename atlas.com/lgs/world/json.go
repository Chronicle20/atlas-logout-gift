package world

import "atlas-lgs/world/gift"

type JSONModel struct {
	Id    int32            `json:"world"`
	Gifts []gift.JSONModel `json:"gifts"`
}
