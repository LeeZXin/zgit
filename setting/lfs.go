package setting

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"io"
	"time"
)

var (
	lfsEnabled = static.GetBool("lfs.enabled")

	jwtAuthExpiry time.Duration

	jwtSecretBytes = make([]byte, 32)

	maxLfsFileSize int64 = 100 * 1024 * 1024
)

func init() {
	jwtExpiry := static.GetInt("lfs.jwt.expiry")
	if jwtExpiry > 0 {
		jwtAuthExpiry = time.Duration(jwtExpiry) * time.Second
	} else {
		jwtAuthExpiry = 24 * time.Hour
	}
	jwtSecretBase64 := static.GetString("lfs.jwt.key")
	if jwtSecretBase64 == "" {
		var err error
		jwtSecretBytes, err = newJwtSecret()
		if err != nil {
			logger.Logger.Panic(err)
		}
	} else {
		n, err := base64.RawURLEncoding.Decode(jwtSecretBytes, []byte(jwtSecretBase64))
		if err != nil || n != 32 {
			jwtSecretBytes, err = newJwtSecret()
			if err != nil {
				logger.Logger.Panic(err)
			}
		}
	}
}

func LfsEnabled() bool {
	return lfsEnabled
}

func LfsJwtAuthExpiry() time.Duration {
	return jwtAuthExpiry
}

func LfsJwtSecretBytes() []byte {
	return jwtSecretBytes
}

// newJwtSecret generates a new value intended to be used for JWT secrets.
func newJwtSecret() ([]byte, error) {
	bytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func MaxLfsFileSize() int64 {
	return maxLfsFileSize
}
