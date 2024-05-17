package user_entity

type Role string

const (
	IT    Role = "it"
	Nurse Role = "nurse"
)

type User struct {
	ID   string `json:"id"`
	NIP  int    `json:"nip"`
	Name string `json:"name"`
}

type RegisterITUser struct {
	NIP      int    `json:"nip" validate:"required,nip"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type RegisterNurseUser struct {
	NIP      int    `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginUser struct {
	NIP      int    `json:"nip" validate:"required,nip"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type LoggedInUser struct {
	ID          string `json:"userId"`
	NIP         int    `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
