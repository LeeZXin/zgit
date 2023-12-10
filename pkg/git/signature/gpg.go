package signature

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"hash"
	"io"
	"strings"
	"time"
)

const (
	StartGPGSigLineTag = "-----BEGIN PGP SIGNATURE-----"
	EndGPGSigLineTag   = "-----END PGP SIGNATURE-----"
)

type GPGSignature struct {
	*packet.Signature
}

func (s *GPGSignature) HashObject(msg []byte) (hash.Hash, error) {
	h := s.Hash.New()
	if _, err := h.Write(msg); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *GPGSignature) GetGPGSignatureKeyId() string {
	if s.IssuerKeyId == nil || *s.IssuerKeyId == 0 {
		return ""
	}
	return fmt.Sprintf("%X", *s.IssuerKeyId)
}

type GPGPublicKey struct {
	*packet.PublicKey
	PrimaryKeyId string
}

type GPGSig string

func (s GPGSig) IsGPGSig() bool {
	return strings.HasPrefix(string(s), StartGPGSigLineTag)
}

func (s GPGSig) IsSSHSig() bool {
	return strings.HasPrefix(string(s), StartSSHSigLineTag)
}

func (s GPGSig) String() string {
	return string(s)
}

func ConvertArmoredGPGKeyString(content string) (openpgp.EntityList, error) {
	list, err := openpgp.ReadArmoredKeyRing(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	return CoalesceGpgEntityList(list), nil
}

func GetGPGEntityListPublicKeys(entityList openpgp.EntityList) []*GPGPublicKey {
	ret := make([]*GPGPublicKey, 0)
	for _, entity := range entityList {
		primaryKeyId := entity.PrimaryKey.KeyIdString()
		ret = append(ret, &GPGPublicKey{
			PublicKey:    entity.PrimaryKey,
			PrimaryKeyId: primaryKeyId,
		})
		for _, subkey := range entity.Subkeys {
			ret = append(ret, &GPGPublicKey{
				PublicKey:    subkey.PublicKey,
				PrimaryKeyId: primaryKeyId,
			})
		}
	}
	return ret
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
	signer, err := openpgp.CheckArmoredDetachedSignature(
		ekeys,
		strings.NewReader(token),
		strings.NewReader(signature),
	)
	if err != nil {
		signer, err = openpgp.CheckArmoredDetachedSignature(
			ekeys,
			strings.NewReader(token+"\n"),
			strings.NewReader(signature),
		)
		if err != nil {
			signer, err = openpgp.CheckArmoredDetachedSignature(
				ekeys,
				strings.NewReader(token+"\r\n"),
				strings.NewReader(signature),
			)
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
				if _, has = original.Identities[name]; has {
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

// GetGPGKeyExpiryTime extract the expiry time of primary key based on sig
func GetGPGKeyExpiryTime(e *openpgp.Entity) time.Time {
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

func ParseGPGSignature(content GPGSig) (*GPGSignature, error) {
	block, err := armor.Decode(strings.NewReader(content.String()))
	if err != nil {
		return nil, err
	}
	if block.Type != openpgp.SignatureType {
		return nil, fmt.Errorf("expected '" + openpgp.SignatureType + "', got: " + block.Type)
	}
	p, err := packet.Read(block.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature packet:%v", err)
	}
	sig, ok := p.(*packet.Signature)
	if !ok {
		return nil, fmt.Errorf("packet is not a signature")
	}
	return &GPGSignature{
		Signature: sig,
	}, nil
}

func ParseGPGPublicKey(content string) (*packet.PublicKey, error) {
	block, err := armor.Decode(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	if block.Type != "PGP PUBLIC KEY BLOCK" {
		return nil, fmt.Errorf("expected '" + openpgp.SignatureType + "', got: " + block.Type)
	}
	p, err := packet.Read(block.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key packet:%v", err)
	}
	sig, ok := p.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("packet is not a public key")
	}
	return sig, nil
}

// Base64EncGPGPubKey encode public key content to base 64
func Base64EncGPGPubKey(pubkey *packet.PublicKey) (string, error) {
	var w bytes.Buffer
	err := pubkey.Serialize(&w)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(w.Bytes()), nil
}

func readerFromBase64(s string) (io.Reader, error) {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(bs), nil
}

// Base64DecGPGPubKey decode public key content from base 64
func Base64DecGPGPubKey(content string) (*packet.PublicKey, error) {
	b, err := readerFromBase64(content)
	if err != nil {
		return nil, err
	}
	// Read key
	p, err := packet.Read(b)
	if err != nil {
		return nil, err
	}
	// Check type
	pkey, ok := p.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not a public key")
	}
	return pkey, nil
}
