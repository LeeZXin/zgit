package usermd

type UserInfo struct {
	UserId       string `json:"userId"`
	Account      string `json:"account"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	CorpId       string `json:"corpId"`
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
	CorpId    string
	AvatarUrl string
}
