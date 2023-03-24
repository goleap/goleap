package models

import "time"

type PostsModel struct {
	Id       uint            `dbKit:"column:id, primaryKey"`
	User     UsersModel      `dbKit:"column:user_id, foreignKey:id"`
	Comments []CommentsModel `dbKit:"column:comments_id, foreignKey:id"`
	Title    string          `dbKit:"column:title"`
	Content  string          `dbKit:"column:content"`
	Created  time.Time       `dbKit:"column:created_at"`
	Updated  time.Time       `dbKit:"column:updated_at"`
}

func (s *PostsModel) DatabaseName() string {
	return "acceptance"
}

func (s *PostsModel) TableName() string {
	return "posts"
}
