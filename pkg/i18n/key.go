package i18n

type Key string

const (
	SystemInternalError Key = "system.internalErr"
	SystemNotAdmin      Key = "system.notAdmin"
	SystemInvalidArgs   Key = "system.invalidArgs"
	SystemNotLogin      Key = "system.notLogin"
	SystemUnauthorized  Key = "system.unauthorized"
	SystemUnimplemented Key = "system.unimplemented"
	SystemAlreadyExists Key = "system.dataAlreadyExists"
)

const (
	UserInvalidId       Key = "user.invalidId"
	UserInvalidName     Key = "user.invalidName"
	UserInvalidEmail    Key = "user.invalidEmail"
	UserInvalidCorpId   Key = "user.invalidCorpId"
	UserInvalidAccount  Key = "user.invalidAccount"
	UserInvalidPassword Key = "user.invalidPassword"
	UserNotFound        Key = "user.notFound"
	UserWrongPassword   Key = "user.wrongPassword"
	UserAlreadyExists   Key = "user.alreadyExists"

	UserAccountNotFoundWarnFormat Key = "user.notFoundWarnFormat"

	UserAccountUnauthorizedReviewCodeWarnFormat Key = "user.unauthorizedReviewCodeWarnFormat"
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
	PrivateRepoType  Key = "repo.privateType"
)

const (
	ProjectInvalidId Key = "project.invalidId"
	ProjectNotFound  Key = "project.notFound"

	ProjectAdminUserGroupName           Key = "project.adminUserGroupName"
	ProjectUserGroupHasUserWhenDel      Key = "project.userGroupHasUserWhenDel"
	ProjectUserGroupUpdateAdminNotAllow Key = "project.updateProjectUserAdminGroupNotAllow"
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

const (
	TimeBeforeSecondUnit Key = "time.beforeSecondUnit"
	TimeBeforeMinuteUnit Key = "time.beforeMinuteUnit"
	TimeBeforeHourUnit   Key = "time.beforeHourUnit"
	TimeBeforeDayUnit    Key = "time.beforeDdayUnit"
	TimeBeforeMonthUnit  Key = "time.beforeMonthUnit"
	TimeBeforeYearUnit   Key = "time.beforeYearUnit"
)

const (
	RepoOpenStatus    Key = "repo.openStatus"
	RepoClosedStatus  Key = "repo.closedStatus"
	RepoDeletedStatus Key = "repo.deletedStatus"
	RepoUnknownStatus Key = "repo.unknownStatus"
)

const (
	PullRequestCannotMerge   Key = "pullRequest.cannotMerge"
	PullRequestOpenStatus    Key = "pullRequest.openStatus"
	PullRequestClosedStatus  Key = "pullRequest.closedStatus"
	PullRequestMergedStatus  Key = "pullRequest.mergedStatus"
	PullRequestUnknownStatus Key = "pullRequest.unknownStatus"
	PullRequestMergeMessage  Key = "pullRequest.mergeMessage"

	PullRequestAgreeMergeStatus    Key = "pullRequest.agreeMerge"
	PullRequestDisagreeMergeStatus Key = "pullRequest.disagreeMerge"
	PullRequestUnknownReviewStatus Key = "pullRequest.unknownReviewStatus"

	PullRequestReviewerCountLowerThanCfg Key = "pullRequest.reviewerCountLowerThanCfg"
)

const (
	SshKeyAlreadyVerified    Key = "sshKey.alreadyVerified"
	SshKeyVerifyTokenExpired Key = "sshKey.verifyTokenExpired"
	SshKeyVerifyFailed       Key = "sshKey.verifyFailed"
)

const (
	ProtectedBranchInvalidReviewCountWhenCreatePr Key = "protectedBranch.invalidReviewCountWhenCreatePr"
	ProtectedBranchNotAllowForcePush              Key = "protectedBranch.notAllowForcePush"
	ProtectedBranchNotAllowDelete                 Key = "protectedBranch.notAllowDelete"
)

var (
	defaultRetMap = map[Key]string{
		SystemInternalError: "系统异常",
		SystemNotAdmin:      "您不是管理员",
		SystemInvalidArgs:   "参数错误",
		SystemNotLogin:      "未登录",
		SystemUnauthorized:  "权限不足",
		SystemUnimplemented: "方法未实现",
		SystemAlreadyExists: "数据已存在",

		UserInvalidId:                 "用户id不合法",
		UserInvalidName:               "用户姓名不合法",
		UserInvalidEmail:              "用户邮箱不合法",
		UserInvalidCorpId:             "企业Id不合法",
		UserInvalidAccount:            "用户账号不合法",
		UserInvalidPassword:           "用户密码不合法",
		UserNotFound:                  "用户不存在",
		UserWrongPassword:             "密码不正确",
		UserAlreadyExists:             "用户已存在",
		UserAccountNotFoundWarnFormat: "用户%s不存在",

		SshKeyFormatError:    "ssh公钥格式错误",
		SshKeyAlreadyExists:  "ssh公钥已存在",
		SshKeyInvalidName:    "ssh公钥名称不合法",
		SshKeyInvalidKeyType: "ssh公钥类型错误",
		SshKeyNotFound:       "ssh公钥不存在",
		SshKeyInvalidKeyId:   "ssh公钥id不合法",

		InternalRepoType: "普通仓库",
		PublicRepoType:   "开源仓库",
		UnKnownRepoType:  "未知类型",
		PrivateRepoType:  "私有仓库",

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

		ProjectAdminUserGroupName:           "项目管理员",
		ProjectUserGroupHasUserWhenDel:      "该项目组存在关联用户, 无法删除",
		ProjectUserGroupUpdateAdminNotAllow: "不允许编辑项目管理员权限",

		TimeBeforeSecondUnit: "秒前",
		TimeBeforeMinuteUnit: "分钟前",
		TimeBeforeHourUnit:   "小时前",
		TimeBeforeDayUnit:    "天前",
		TimeBeforeMonthUnit:  "月前",
		TimeBeforeYearUnit:   "年前",

		RepoOpenStatus:    "打开状态",
		RepoClosedStatus:  "归档状态",
		RepoUnknownStatus: "未知状态",
		RepoDeletedStatus: "删除状态",

		PullRequestCannotMerge:   "无法合并",
		PullRequestOpenStatus:    "已打开",
		PullRequestClosedStatus:  "已关闭",
		PullRequestMergedStatus:  "已合并",
		PullRequestUnknownStatus: "未知",
		PullRequestMergeMessage:  "合并请求: %s, 申请人: %s, 合并人: %s",

		SshKeyAlreadyVerified:    "已经校验过",
		SshKeyVerifyTokenExpired: "token已失效",
		SshKeyVerifyFailed:       "校验失败",

		ProtectedBranchInvalidReviewCountWhenCreatePr: "保护分支代码评审者数量不合法",

		UserAccountUnauthorizedReviewCodeWarnFormat: "该用户%s无评审代码的权限",

		ProtectedBranchNotAllowForcePush: "保护分支禁止强制push",
		ProtectedBranchNotAllowDelete:    "保护分支禁止删除",

		PullRequestAgreeMergeStatus:    "同意合并",
		PullRequestDisagreeMergeStatus: "不同意合并",
		PullRequestUnknownReviewStatus: "未知状态",

		PullRequestReviewerCountLowerThanCfg: "代码评审数量小于配置数量",
	}
)

func (k Key) DefaultRet() string {
	return defaultRetMap[k]
}
