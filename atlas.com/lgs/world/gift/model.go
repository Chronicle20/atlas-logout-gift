package gift

type Model struct {
	id           uint32
	worldId      int32
	itemId       uint32
	serialNumber uint32
	weight       uint32
}

func (m Model) Weight() uint32 {
	return m.weight
}

func (m Model) SerialNumber() uint32 {
	return m.serialNumber
}

type builder struct {
	id           uint32
	worldId      int32
	itemId       uint32
	serialNumber uint32
	weight       uint32
}

func NewGiftBuilder(id uint32) *builder {
	return &builder{id: id}
}

func (b *builder) SetWorldId(worldId int32) *builder {
	b.worldId = worldId
	return b
}

func (b *builder) SetItemId(itemId uint32) *builder {
	b.itemId = itemId
	return b
}

func (b *builder) SetSerialNumber(serialNumber uint32) *builder {
	b.serialNumber = serialNumber
	return b
}

func (b *builder) SetWeight(weight uint32) *builder {
	b.weight = weight
	return b
}

func (b *builder) Build() Model {
	return Model{
		id:           b.id,
		worldId:      b.worldId,
		itemId:       b.itemId,
		serialNumber: b.serialNumber,
		weight:       b.weight,
	}
}
