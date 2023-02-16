package helper

type ExtraModel struct {
	Id        uint       `goleap:"column:id"`
	BaseModel *BaseModel `goleap:"column:recursive_id, foreignKey:id"`
	//Slice     []BaseModel    `goleap:"column:slice_id, foreignKey:id"`
	ExtraJump ExtraJumpModel `goleap:"column:extra_jump_id, foreignKey:id"`
}

func (s ExtraModel) DatabaseName() string {
	return "test"
}

func (s ExtraModel) TableName() string {
	return "extra"
}
