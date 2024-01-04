package signature

import (
	"bytes"
	"github.com/42wim/sshsig"
)

const (
	StartSSHSigLineTag = "-----BEGIN SSH SIGNATURE-----"
	EndSSHSigLineTag   = "-----END SSH SIGNATURE-----"
	DefaultNamespace   = "git"
)

func VerifySshSignature(sig, payload, publicKey string) error {
	return sshsig.Verify(bytes.NewBuffer([]byte(payload)), []byte(sig), []byte(publicKey), DefaultNamespace)
}
