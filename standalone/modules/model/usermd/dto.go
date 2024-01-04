package usermd

type UserInfo struct {
	Account      string `json:"account"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsAdmin      bool   `json:"isAdmin"`
	IsProhibited bool   `json:"isProhibited"`
	AvatarUrl    string `json:"avatarUrl"`
}

type InsertUserReqDTO struct {
	Account   string
	Name      string
	Email     string
	Password  string
	IsAdmin   bool
	AvatarUrl string
}

type ListUserReqDTO struct {
	Account string
	Offset  int64
	Limit   int
}

type UpdateUserReqDTO struct {
	Account string
	Name    string
	Email   string
}

type UpdateAdminReqDTO struct {
	Account string
	IsAdmin bool
}

type UpdatePasswordReqDTO struct {
	Account  string
	Password string
}
