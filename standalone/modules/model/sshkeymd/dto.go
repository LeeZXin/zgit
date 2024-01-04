package sshkeymd

type InsertSshKeyReqDTO struct {
	Account     string
	Name        string
	Fingerprint string
	Content     string
	Verified    bool
}

type ListSshKeyReqDTO struct {
	Offset  int64
	Limit   int
	Account string
}

type UpdateVerifiedVarReqDTO struct {
	KeyId    string
	Verified bool
}
