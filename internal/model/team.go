package model

type Team struct {
	Id     int    `json:"id,omitempty" db:"id"`
	UserId int    `json:"user_id,omitempty" db:"user_id"`
	Title  string `json:"title" db:"title" binding:"required,min=1"`
}

type UpdateTeamInput struct {
	Title string `json:"title" binding:"required,min=3"`
}
