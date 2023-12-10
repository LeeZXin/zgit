package signature

import (
	"bytes"
	"github.com/42wim/sshsig"
	"io"
)

const (
	StartSSHSigLineTag = "-----BEGIN SSH SIGNATURE-----"
	EndSSHSigLineTag   = "-----END SSH SIGNATURE-----"
	DefaultNamespace   = "git"
)

func VerifySshSignature(sig, payload, publicKey string) error {
	return sshsig.Verify(bytes.NewBuffer([]byte(payload)), []byte(sig), []byte(publicKey), DefaultNamespace)
}

func SignSshSignature(privateKey string, data io.Reader) (string, error) {
	sign, err := sshsig.Sign([]byte(privateKey), data, DefaultNamespace)
	if err != nil {
		return "", err
	}
	return string(sign), nil
}
