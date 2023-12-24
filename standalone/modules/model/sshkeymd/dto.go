package sshkeymd

type InsertSshKeyReqDTO struct {
	Account     string `json:"account"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Content     string `json:"content"`
	Verified    bool   `json:"verified"`
}
