package models

type UserModel struct {
	Id uint `dbKit:"column:id, primaryKey"`
}

func (s *UserModel) DatabaseName() string {
	return "acceptance"
}

func (s *UserModel) TableName() string {
	return "user"
}
