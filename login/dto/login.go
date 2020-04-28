package dto

type LoginDTO struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

/*
{
    "app_id": 12345678,
    "app_key": "eyJleHAiOjE1NDc1NzEzNzEsInN1YiI6M30"
}

*/
//
//type OpenAuthDTO struct {
//	Email    string `json:"app_id"`
//	Password string `json:"app_key"`
//}

type ChangePassword struct {
	Token           string
	UserId          int64
	FirstAction     int    `json:"action"`
	Password        string `json:"old_password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,gte=6,lte=40"`
	NewPassword     string `json:"new_password" validate:"required,gte=6,lte=40"`
}
