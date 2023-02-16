package helper

type BaseModel struct {
	Id        uint        `goleap:"column:id, primaryKey"`
	Recursive *BaseModel  `goleap:"column:recursive_id, foreignKey:id"`
	Slice     []BaseModel `goleap:"column:slice_id, foreignKey:id"`
	Extra     ExtraModel  `goleap:"column:extra_id, foreignKey:id"`

	private bool
}

func (s BaseModel) DatabaseName() string {
	return "test"
}

func (s BaseModel) TableName() string {
	return "base"
}
