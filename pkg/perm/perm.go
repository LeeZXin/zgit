package perm

var (
	DefaultProjectPerm = ProjectPerm{
		CanInitRepo:   true,
		CanDeleteRepo: true,
	}
	DefaultRepoPerm = RepoPerm{
		CanAccess: true,
		CanPush:   true,
		CanClose:  true,

		CanHandleProtectedBranch: true,
		CanHandlePullRequest:     true,
	}
	DefaultPermDetail = Detail{
		ProjectPerm:          DefaultProjectPerm,
		ApplyDefaultRepoPerm: true,
		DefaultRepoPerm:      DefaultRepoPerm,
	}
)

type Detail struct {
	// 项目权限
	ProjectPerm ProjectPerm `json:"projectPerm"`
	// 使用仓库全局默认权限
	ApplyDefaultRepoPerm bool `json:"applyDefaultRepoPerm"`
	// 默认仓库权限
	DefaultRepoPerm RepoPerm `json:"defaultRepoPerm"`
	// 可访问仓库权限列表
	RepoPermList []RepoPermWithId `json:"repoPermList,omitempty"`
}

func (d *Detail) GetRepoPerm(repoId string) RepoPerm {
	if d.ApplyDefaultRepoPerm {
		return d.DefaultRepoPerm
	}
	for _, perm := range d.RepoPermList {
		if perm.RepoId == repoId {
			return perm.RepoPerm
		}
	}
	return RepoPerm{}
}

type RepoPermWithId struct {
	RepoId string `json:"repoId"`
	RepoPerm
}

type RepoPerm struct {
	// 可访问
	CanAccess bool `json:"canAccess"`
	// 可推送代码
	CanPush bool `json:"canPush"`
	// 是否可归档
	CanClose bool `json:"canClose"`
	// 是否可处理保护分支
	CanHandleProtectedBranch bool `json:"canHandleProtectedBranch"`
	// 是否可处理pr
	CanHandlePullRequest bool `json:"canHandlePullRequest"`
}

type ProjectPerm struct {
	// 是否可创建仓库
	CanInitRepo bool `json:"canInitRepo"`
	// 是否可删除仓库
	CanDeleteRepo bool `json:"canDeleteRepo"`
}
