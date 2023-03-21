package models

type ExtraJumpModel struct {
	Id         uint        `dbKit:"column:id"`
	JumpToBase *BaseModel  `dbKit:"column:recursive_id, foreignKey:id"`
	Slice      []BaseModel `dbKit:"column:slice_id, foreignKey:id"`
}

func (s *ExtraJumpModel) DatabaseName() string {
	return "test"
}

func (s *ExtraJumpModel) TableName() string {
	return "extra_jump"
}
