package sshkeyapi

type InsertSshPubKeyReqVO struct {
	Name          string `json:"name"`
	PubKeyContent string `json:"pubKeyContent"`
}

type DeleteSshKeyReqVO struct {
	KeyId string `json:"keyId"`
}
