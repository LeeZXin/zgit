package sshkeymd

type InsertSshKeyReqDTO struct {
	UserId      string `json:"userId"`
	CorpId      string `json:"corpId"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Content     string `json:"content"`
	KeyType     int    `json:"keyType"`
	Verified    bool   `json:"verified"`
}
