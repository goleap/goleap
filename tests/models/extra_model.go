package models

type ExtraModel struct {
	BaseModel *BaseModel     `dbKit:"column:recursive_id, foreignKey:id"`
	ExtraJump ExtraJumpModel `dbKit:"column:extra_jump_id, foreignKey:id"`
	Id        uint           `dbKit:"column:id, primaryKey"`
}

func (s ExtraModel) DatabaseName() string {
	return "test"
}

func (s ExtraModel) TableName() string {
	return "extra"
}
