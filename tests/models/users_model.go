package models

import "time"

type UsersModel struct {
	Id        uint       `dbKit:"column:id, primaryKey"`
	Email     string     `dbKit:"column:email"`
	Password  string     `dbKit:"column:password"`
	Validated bool       `dbKit:"column:validated"`
	CreatedAt *time.Time `dbKit:"column:created_at"`
	UpdatedAt *time.Time `dbKit:"column:updated_at"`

	BadTag  string `dbKit:"column:bad_tag:bad_tag, "`
	private bool
}

func (s *UsersModel) DatabaseName() string {

	// This is a test for static-check validation.
	s.private = false

	return "acceptance"
}

func (s *UsersModel) TableName() string {
	return "users"
}
