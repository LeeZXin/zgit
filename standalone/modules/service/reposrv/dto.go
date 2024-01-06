package reposrv

import (
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"regexp"
	"strings"
	"time"
	"zgit/pkg/git"
	"zgit/standalone/modules/model/projectmd"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

const (
	UpDirection   = "up"
	DownDirection = "down"
)

var (
	validRepoNamePattern = regexp.MustCompile("^[\\w\\-]{1,32}$")
	validBranchPattern   = regexp.MustCompile("^\\w{1,32}$")
)

type InitRepoReqDTO struct {
	Operator      usermd.UserInfo
	ProjectId     string
	Name          string
	Desc          string
	RepoType      repomd.RepoType
	CreateReadme  bool
	GitIgnoreName string
	DefaultBranch string
}

func (r *InitRepoReqDTO) IsValid() error {
	if !projectmd.IsProjectIdValid(r.ProjectId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !validRepoNamePattern.MatchString(r.Name) {
		return util.InvalidArgsError()
	}
	if len(r.Desc) > 255 {
		return util.InvalidArgsError()
	}
	if r.DefaultBranch != "" && !validBranchPattern.MatchString(r.DefaultBranch) {
		return util.InvalidArgsError()
	}
	if !r.RepoType.IsValid() {
		return util.InvalidArgsError()
	}
	if r.GitIgnoreName != "" && !gitignoreSet.Contains(r.GitIgnoreName) {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteRepoReqDTO struct {
	RepoId   string
	Operator usermd.UserInfo
}

func (r *DeleteRepoReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	return nil
}

type TreeRepoReqDTO struct {
	RepoId   string
	RefName  string
	Dir      string
	Operator usermd.UserInfo
}

func (r *TreeRepoReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	if len(r.RefName) > 128 || len(r.RefName) == 0 {
		return util.InvalidArgsError()
	}
	if strings.HasSuffix(r.Dir, "/") {
		return util.InvalidArgsError()
	}
	return nil
}

type CatFileReqDTO struct {
	RepoId   string
	RefName  string
	Dir      string
	FileName string
	Operator usermd.UserInfo
}

func (r *CatFileReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	if len(r.RefName) > 128 || len(r.RefName) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.Dir) > 128 || len(r.Dir) == 0 {
		return util.InvalidArgsError()
	}
	if !validateFileName(r.FileName) {
		return util.InvalidArgsError()
	}
	if strings.HasSuffix(r.Dir, "/") {
		return util.InvalidArgsError()
	}
	return nil
}

type EntriesRepoReqDTO struct {
	RepoId   string
	RefName  string
	Dir      string
	Offset   int
	Operator usermd.UserInfo
}

func (r *EntriesRepoReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	if r.RefName == "" {
		return util.InvalidArgsError()
	}
	if strings.HasSuffix(r.Dir, "/") {
		return util.InvalidArgsError()
	}
	if r.Offset < 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type ListRepoReqDTO struct {
	ProjectId string
	Operator  usermd.UserInfo
}

func (r *ListRepoReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if len(r.ProjectId) == 0 || len(r.ProjectId) > 128 {
		return util.InvalidArgsError()
	}
	return nil
}

type CommitDTO struct {
	Author        git.User
	Committer     git.User
	AuthoredDate  time.Time
	CommittedDate time.Time
	CommitMsg     string
	CommitId      string
	ShortId       string
}

type FileDTO struct {
	Mode    string
	RawPath string
	Path    string
	Commit  CommitDTO
}

type TreeDTO struct {
	Files   []FileDTO
	Limit   int
	Offset  int
	HasMore bool
}

type TreeRepoRespDTO struct {
	IsEmpty      bool
	ReadmeText   string
	HasReadme    bool
	RecentCommit CommitDTO
	Tree         TreeDTO
}

type RepoTypeDTO struct {
	Option int
	Name   string
}

type AllBranchesReqDTO struct {
	RepoId   string
	Operator usermd.UserInfo
}

func (r *AllBranchesReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	return nil
}

type AllTagsReqDTO struct {
	RepoId   string
	Operator usermd.UserInfo
}

func (r *AllTagsReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	return nil
}

type GcReqDTO struct {
	RepoId   string
	Operator usermd.UserInfo
}

func (r *GcReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	return nil
}

var gitignoreSet = hashset.NewHashSet([]string{
	"AL", "Actionscript", "Ada", "Agda", "AltiumDesigner", "Android", "Anjuta", "Ansible", "AppEngine",
	"AppceleratorTitanium", "ArchLinuxPackages", "Archives", "AtmelStudio", "AutoIt", "Autotools", "B4X", "Backup",
	"Bazaar", "Bazel", "Beef", "Bitrix", "BricxCC", "C", "C++", "CDK", "CFWheels", "CMake", "CUDA", "CakePHP",
	"Calabash", "ChefCookbook", "Clojure", "Cloud9", "CodeIgniter", "CodeKit", "CodeSniffer", "CommonLisp",
	"Composer", "Concrete5", "Coq", "Cordova", "CraftCMS", "D", "DM", "Dart", "DartEditor", "Delphi", "Diff",
	"Dreamweaver", "Dropbox", "Drupal", "Drupal7", "EPiServer", "Eagle", "Eclipse", "EiffelStudio", "Elisp",
	"Elixir", "Elm", "Emacs", "Ensime", "Erlang", "Espresso", "Exercism", "ExpressionEngine", "ExtJs", "Fancy",
	"Finale", "FlaxEngine", "FlexBuilder", "ForceDotCom", "Fortran", "FuelPHP", "GNOMEShellExtension", "GPG",
	"GWT", "Gcov", "GitBook", "Go", "Go.AllowList", "Godot", "Gradle", "Grails", "Gretl", "Haskell", "Hugo",
	"IAR_EWARM", "IGORPro", "Idris", "Images", "InforCMS", "JBoss", "JBoss4", "JBoss6", "JDeveloper", "JENKINS_HOME",
	"JEnv", "Java", "Jekyll", "JetBrains", "Jigsaw", "Joomla", "Julia", "JupyterNotebooks", "KDevelop4", "Kate",
	"Kentico", "KiCad", "Kohana", "Kotlin", "LabVIEW", "Laravel", "Lazarus", "Leiningen", "LemonStand", "LensStudio",
	"LibreOffice", "Lilypond", "Linux", "Lithium", "Logtalk", "Lua", "LyX", "MATLAB", "Magento", "Magento1", "Magento2",
	"Maven", "Mercurial", "Mercury", "MetaProgrammingSystem", "Metals", "Meteor", "MicrosoftOffice", "ModelSim",
	"Momentics", "MonoDevelop", "NWjs", "Nanoc", "NasaSpecsIntact", "NetBeans", "Nikola", "Nim", "Ninja", "Nix",
	"Node", "NotepadPP", "OCaml", "Objective-C", "Octave", "Opa", "OpenCart", "OpenSSL", "OracleForms", "Otto",
	"PSoCCreator", "Packer", "Patch", "Perl", "Perl6", "Phalcon", "Phoenix", "Pimcore", "PlayFramework", "Plone",
	"Prestashop", "Processing", "PuTTY", "Puppet", "PureScript", "Python", "Qooxdoo", "Qt", "R", "ROS", "ROS2",
	"Racket", "Rails", "Raku", "Red", "Redcar", "Redis", "RhodesRhomobile", "Ruby", "Rust", "SAM", "SBT", "SCons",
	"SPFx", "SVN", "Sass", "Scala", "Scheme", "Scrivener", "Sdcc", "SeamGen", "SketchUp", "SlickEdit", "Smalltalk",
	"Snap", "Splunk", "Stata", "Stella", "Strapi", "SublimeText", "SugarCRM", "Swift", "Symfony", "SymphonyCMS",
	"Syncthing", "SynopsysVCS", "Tags", "TeX", "Terraform", "TextMate", "Textpattern", "ThinkPHP", "Toit", "TortoiseGit",
	"TurboGears2", "TwinCAT3", "Typo3", "Umbraco", "Unity", "UnrealEngine", "V", "VVVV", "Vagrant", "Vim", "VirtualEnv",
	"Virtuoso", "VisualStudio", "VisualStudioCode", "Vue", "Waf", "WebMethods", "Windows", "WordPress", "Xcode", "Xilinx",
	"XilinxISE", "Xojo", "Yeoman", "Yii", "ZendFramework", "Zephir", "core", "esp-idf", "macOS", "uVision",
})

type DiffCommitsReqDTO struct {
	RepoId   string
	Target   string
	Head     string
	Operator usermd.UserInfo
}

func (r *DiffCommitsReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Target) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Head) {
		return util.InvalidArgsError()
	}
	return nil
}

type DiffCommitsRespDTO struct {
	Target        string              `json:"target"`
	Head          string              `json:"head"`
	TargetCommit  CommitDTO           `json:"targetCommit"`
	HeadCommit    CommitDTO           `json:"headCommit"`
	Commits       []CommitDTO         `json:"commits"`
	NumFiles      int                 `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoDTO `json:"diffNumsStats"`
	ConflictFiles []string            `json:"conflictFiles"`
	CanMerge      bool                `json:"canMerge"`
}

type DiffNumsStatInfoDTO struct {
	FileChangeNums int `json:"fileChangeNums"`
	InsertNums     int `json:"insertNums"`
	DeleteNums     int `json:"deleteNums"`
	Stats          []DiffNumsStatDTO
}

type DiffNumsStatDTO struct {
	RawPath    string
	Path       string
	TotalNums  int
	InsertNums int
	DeleteNums int
}

type DiffFileRespDTO struct {
	FilePath    string
	OldMode     string
	Mode        string
	IsSubModule bool
	FileType    git.DiffFileType
	IsBinary    bool
	RenameFrom  string
	RenameTo    string
	CopyFrom    string
	CopyTo      string
	Lines       []DiffLineDTO
}

type DiffLineDTO struct {
	Index   int
	LeftNo  int
	Prefix  string
	RightNo int
	Text    string
}

type ShowDiffTextContentReqDTO struct {
	RepoId    string
	CommitId  string
	FileName  string
	Offset    int
	Limit     int
	Direction string
	Operator  usermd.UserInfo
}

func (r *ShowDiffTextContentReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateCommitId(r.CommitId) {
		return util.InvalidArgsError()
	}
	if r.Offset < 0 {
		return util.InvalidArgsError()
	}
	if len(r.Direction) == 0 || len(r.Direction) > 10 {
		return util.InvalidArgsError()
	}
	if r.Direction != UpDirection && r.Direction != DownDirection {
		return util.InvalidArgsError()
	}
	if !validateFileName(r.FileName) {
		return util.InvalidArgsError()
	}
	return nil
}

type DiffFileReqDTO struct {
	RepoId   string
	Target   string
	Head     string
	FileName string
	Operator usermd.UserInfo
}

func (r *DiffFileReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoIdValid(r.RepoId) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Target) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Head) {
		return util.InvalidArgsError()
	}
	if !validateFileName(r.FileName) {
		return util.InvalidArgsError()
	}
	return nil
}

func validateFileName(name string) bool {
	return len(name) <= 255 && len(name) > 0
}
