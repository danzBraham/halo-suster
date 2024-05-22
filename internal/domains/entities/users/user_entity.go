package user_entity

import "time"

type Role string

const (
	IT    Role = "it"
	Nurse Role = "nurse"
)

type User struct {
	ID        string    `json:"id"`
	NIP       int       `json:"nip"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterITUser struct {
	NIP      int    `json:"nip" validate:"required,nip"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type RegisterNurseUser struct {
	NIP          int    `json:"nip" validate:"required,nip"`
	Name         string `json:"name" validate:"required,min=5,max=50"`
	CardImageURL string `json:"identityCardScanImg" validate:"required,imageurl"`
}

type LoginUser struct {
	NIP      int    `json:"nip" validate:"required,nip"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type LoggedInUser struct {
	UserID      string `json:"userId"`
	NIP         int    `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type UserQueryParams struct {
	UserID    string
	Limit     int
	Offset    int
	Name      string
	NIP       string
	Role      string
	CreatedAt string
}

type UserList struct {
	ID        string    `json:"userId"`
	NIP       int       `json:"nip"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateNurseUser struct {
	UserID string `json:"userId"`
	NIP    int    `json:"nip" validate:"required,nip"`
	Name   string `json:"name" validate:"required,min=5,max=50"`
}

type GiveAccessNurseUser struct {
	UserID   string `json:"userId"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}
