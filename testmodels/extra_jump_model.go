package testmodels

type ExtraJumpModel struct {
	Id         uint        `goleap:"column:id"`
	JumpToBase *BaseModel  `goleap:"column:recursive_id, foreignKey:id"`
	Slice      []BaseModel `goleap:"column:slice_id, foreignKey:id"`
}

func (s ExtraJumpModel) DatabaseName() string {
	return "test"
}

func (s ExtraJumpModel) TableName() string {
	return "extra_jump"
}
