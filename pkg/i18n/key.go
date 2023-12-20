package i18n

type Key string

const (
	SystemInternalError Key = "system.internalErr"
	SystemNotAdmin      Key = "system.notAdmin"
	SystemInvalidArgs   Key = "system.invalidArgs"
	SystemNotLogin      Key = "system.notLogin"
	SystemUnauthorized  Key = "system.unauthorized"
)

const (
	UserInvalidId        Key = "user.invalidId"
	UserInvalidName      Key = "user.invalidName"
	UserInvalidEmail     Key = "user.invalidEmail"
	UserInvalidCorpId    Key = "user.invalidCorpId"
	UserInvalidAccount   Key = "user.invalidAccount"
	UserInvalidSessionId Key = "user.invalidSessionId"
	UserInvalidPassword  Key = "user.invalidPassword"
	UserNotFound         Key = "user.notFound"
	UserWrongPassword    Key = "user.wrongPassword"
)

const (
	SshKeyFormatError    Key = "sshKey.formatErr"
	SshKeyAlreadyExists  Key = "sshKey.alreadyExists"
	SshKeyNotFound       Key = "sshKey.notFound"
	SshKeyInvalidName    Key = "sshKey.invalidName"
	SshKeyInvalidKeyType Key = "sshKey.invalidKeyType"
	SshKeyInvalidKeyId   Key = "sshKey.invalidKeyId"
)

const (
	InternalRepoType Key = "repo.internalType"
	PublicRepoType   Key = "repo.publicType"
	UnKnownRepoType  Key = "repo.unknownType"
)

const (
	ProjectInvalidId Key = "project.invalidId"
	ProjectNotFound  Key = "project.notFound"
)

const (
	RepoInvalidName       Key = "repo.invalidName"
	RepoInvalidDescLength Key = "repo.invalidDescLength"

	RepoInvalidBranch Key = "repo.invalidBranch"
	RepoInitFail      Key = "repo.initFail"
	RepoNotFound      Key = "repo.notFound"
	RepoAlreadyExists Key = "repo.alreadyExists"

	RepoInvalidId            Key = "repo.invalidId"
	RepoInvalidType          Key = "repo.invalidType"
	RepoInvalidGitIgnoreName Key = "repo.invalidGitIgnore"

	RepoCountOutOfLimit Key = "repo.countOutOfLimit"
)

const (
	CorpEmptyId Key = "corp.emptyId"
)

var (
	defaultRetMap = map[Key]string{
		SystemInternalError: "系统异常",
		SystemNotAdmin:      "您不是管理员",
		SystemInvalidArgs:   "参数错误",
		SystemNotLogin:      "未登录",
		SystemUnauthorized:  "权限不足",

		UserInvalidId:       "用户id不合法",
		UserInvalidName:     "用户姓名不合法",
		UserInvalidEmail:    "用户邮箱不合法",
		UserInvalidCorpId:   "企业Id不合法",
		UserInvalidAccount:  "用户账号不合法",
		UserInvalidPassword: "用户密码不合法",
		UserNotFound:        "用户不存在",
		UserWrongPassword:   "密码不正确",

		SshKeyFormatError:    "ssh公钥格式错误",
		SshKeyAlreadyExists:  "ssh公钥已存在",
		SshKeyInvalidName:    "ssh公钥名称不合法",
		SshKeyInvalidKeyType: "ssh公钥类型错误",
		SshKeyNotFound:       "ssh公钥不存在",
		SshKeyInvalidKeyId:   "ssh公钥id不合法",

		InternalRepoType: "普通仓库",
		PublicRepoType:   "开源仓库",
		UnKnownRepoType:  "未知类型",

		RepoInvalidName:          "仓库名称不合法",
		RepoInvalidDescLength:    "仓库描述长度不合法",
		RepoInvalidBranch:        "仓库分支不合法",
		RepoInitFail:             "仓库初始化失败",
		RepoAlreadyExists:        "仓库已存在",
		RepoInvalidType:          "仓库类型错误",
		RepoInvalidGitIgnoreName: "gitIgnore名称错误",
		RepoNotFound:             "仓库不存在",
		RepoCountOutOfLimit:      "仓库数量大于上限",
		RepoInvalidId:            "仓库id不合法",

		CorpEmptyId: "公司id为空",

		ProjectInvalidId: "项目id不合法",
		ProjectNotFound:  "项目不存在",
	}
)

func (k Key) DefaultRet() string {
	return defaultRetMap[k]
}
