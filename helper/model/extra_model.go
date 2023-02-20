package model

type ExtraModel struct {
	BaseModel *BaseModel     `goleap:"column:recursive_id, foreignKey:id"`
	ExtraJump ExtraJumpModel `goleap:"column:extra_jump_id, foreignKey:id"`
	Id        uint           `goleap:"column:id, primaryKey"`
}

func (s ExtraModel) DatabaseName() string {
	return "test"
}

func (s ExtraModel) TableName() string {
	return "extra"
}
