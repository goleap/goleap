package model

type UserModel struct {
	Id uint `goleap:"column:id, primaryKey"`
}

func (s *UserModel) DatabaseName() string {
	return "acceptance"
}

func (s *UserModel) TableName() string {
	return "user"
}
