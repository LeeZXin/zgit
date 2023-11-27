package gpg

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"strings"
	"time"
)

const (
	StartLineTag = "-----BEGIN PGP SIGNATURE-----"
	EndLineTag   = "-----END PGP SIGNATURE-----"
)

func ConvertArmoredGPGKeyString(content string) (openpgp.EntityList, error) {
	return openpgp.ReadArmoredKeyRing(strings.NewReader(content))
}

func GetVerifyToken(user User, minutes int) string {
	h := sha256.New()
	_, _ = h.Write([]byte(strings.Join(
		[]string{
			time.Now().Truncate(time.Minute).Add(time.Duration(minutes) * time.Minute).Format(time.RFC1123Z),
			user.CreatedUnix.Format(time.RFC1123Z),
			user.Name,
			user.Email,
			user.Id,
		}, ":",
	)))
	return hex.EncodeToString(h.Sum(nil))
}

func CheckArmoredDetachedSignature(ekeys openpgp.EntityList, token, signature string) (*openpgp.Entity, error) {
	signer, err := openpgp.CheckArmoredDetachedSignature(ekeys, strings.NewReader(token), strings.NewReader(signature))
	if err != nil {
		signer, err = openpgp.CheckArmoredDetachedSignature(ekeys, strings.NewReader(token+"\n"), strings.NewReader(signature))
		if err != nil {
			signer, err = openpgp.CheckArmoredDetachedSignature(ekeys, strings.NewReader(token+"\r\n"), strings.NewReader(signature))
		}
	}
	return signer, err
}

func CoalesceGpgEntityList(ekeys openpgp.EntityList) openpgp.EntityList {
	id2key := map[string]*openpgp.Entity{}
	newEKeys := make([]*openpgp.Entity, 0, len(ekeys))
	for _, ekey := range ekeys {
		id := ekey.PrimaryKey.KeyIdString()
		if original, has := id2key[id]; has {
			// Coalesce this with the other one
			for _, subKey := range ekey.Subkeys {
				if subKey.PublicKey == nil {
					continue
				}
				found := false
				for _, originalSubKey := range original.Subkeys {
					if originalSubKey.PublicKey == nil {
						continue
					}
					if originalSubKey.PublicKey.KeyId == subKey.PublicKey.KeyId {
						found = true
						break
					}
				}
				if !found {
					original.Subkeys = append(original.Subkeys, subKey)
				}
			}
			for name, identity := range ekey.Identities {
				if _, has := original.Identities[name]; has {
					continue
				}
				original.Identities[name] = identity
			}
			continue
		}
		id2key[id] = ekey
		newEKeys = append(newEKeys, ekey)
	}
	return newEKeys
}

// GetExpiryTime extract the expiry time of primary key based on sig
func GetExpiryTime(e *openpgp.Entity) time.Time {
	expiry := time.Time{}
	// Extract self-sign for expire date based on : https://github.com/golang/crypto/blob/master/openpgp/keys.go#L165
	var selfSig *packet.Signature
	for _, ident := range e.Identities {
		if selfSig == nil {
			selfSig = ident.SelfSignature
		} else if ident.SelfSignature.IsPrimaryId != nil && *ident.SelfSignature.IsPrimaryId {
			selfSig = ident.SelfSignature
			break
		}
	}
	if selfSig.KeyLifetimeSecs != nil {
		expiry = e.PrimaryKey.CreationTime.Add(time.Duration(*selfSig.KeyLifetimeSecs) * time.Second)
	}
	return expiry
}
