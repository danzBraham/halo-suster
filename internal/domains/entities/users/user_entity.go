package user_entity

type Role string

const (
	IT    Role = "it"
	Nurse Role = "nurse"
)

type RegisterITUser struct {
	Nip      int    `json:"nip" validate:"required,nip"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type RegisterNurseUser struct {
	Nip      int    `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginUser struct {
	Nip      int    `json:"nip"`
	Password string `json:"password"`
}

type LoggedInUser struct {
	ID          string `json:"userId"`
	Nip         int    `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
