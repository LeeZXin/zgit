package gitsrv

import "zgit/pkg/perm"

const (
	lfsAuthenticateVerb = "git-lfs-authenticate"
)

var (
	hiWords = "Hi there! You've successfully authenticated with the deploy key named %v, but zgit does not provide shell access."

	allowedCommands = map[string]perm.AccessMode{
		"git-upload-pack":    perm.AccessModeRead,
		"git-upload-archive": perm.AccessModeRead,
		"git-receive-pack":   perm.AccessModeWrite,
		lfsAuthenticateVerb:  perm.AccessModeNone,
	}
)

type ServCommandResults struct {
	IsWiki    bool
	RepoId    string
	ClusterId string
}
