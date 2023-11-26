package git

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LeeZXin/zsf-utils/localcache"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"zgit/git/command"
)

var (
	VersionRegexp = regexp.MustCompile(`^(?:\w+\-)*[vV]?` +
		`([0-9]+(\.[0-9]+)*?)` +
		`(-` +
		`([0-9]+[0-9A-Za-z\-~]*(\.[0-9A-Za-z\-~]+)*)` +
		`|` +
		`([-\.]?([A-Za-z\-~]+[0-9A-Za-z\-~]*(\.[0-9A-Za-z\-~]+)*)))?` +
		`(\+([0-9A-Za-z\-~]+(\.[0-9A-Za-z\-~]+)*))?` +
		`([\+\.\-~]g[0-9A-Fa-f]{10}$)?` +
		`?$`)
	versionCache, _ = localcache.NewLazyLoader[*Version](getGitVersion)
)

type Version struct {
	metadata string
	pre      string
	segments []int64
	si       int
	original string
}

func newVersion(v string) (*Version, error) {
	matches := VersionRegexp.FindStringSubmatch(v)
	if matches == nil {
		return nil, fmt.Errorf("malformed version: %s", v)
	}
	segmentsStr := strings.Split(matches[1], ".")
	segments := make([]int64, len(segmentsStr))
	si := 0
	for i, str := range segmentsStr {
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %s", err)
		}
		segments[i] = val
		si++
	}
	for i := len(segments); i < 3; i++ {
		segments = append(segments, 0)
	}
	pre := matches[7]
	if pre == "" {
		pre = matches[4]
	}
	return &Version{
		metadata: matches[10],
		pre:      pre,
		segments: segments,
		si:       si,
		original: v,
	}, nil
}

func (v *Version) String() string {
	var buf bytes.Buffer
	fmtParts := make([]string, len(v.segments))
	for i, s := range v.segments {
		str := strconv.FormatInt(s, 10)
		fmtParts[i] = str
	}
	fmt.Fprint(&buf, strings.Join(fmtParts, "."))
	if v.pre != "" {
		fmt.Fprintf(&buf, "-%s", v.pre)
	}
	if v.metadata != "" {
		fmt.Fprintf(&buf, "+%s", v.metadata)
	}
	return buf.String()
}

func (v *Version) copySegments() []int64 {
	ret := make([]int64, len(v.segments))
	copy(ret, v.segments)
	return ret
}

func (v *Version) Original() string {
	return v.original
}

func (v *Version) Compare(other *Version) int {
	if v.String() == other.String() {
		return 0
	}
	ss := v.copySegments()
	so := other.copySegments()
	if reflect.DeepEqual(ss, so) {
		ps := v.pre
		po := other.pre
		if ps == "" && po == "" {
			return 0
		}
		if ps == "" {
			return 1
		}
		if po == "" {
			return -1
		}
		return comparePre(ps, po)
	}
	ls := len(ss)
	lo := len(so)
	hs := ls
	if ls < lo {
		hs = lo
	}
	for i := 0; i < hs; i++ {
		if i > ls-1 {
			if !allZero(so[i:]) {
				return -1
			}
			break
		} else if i > lo-1 {
			if !allZero(ss[i:]) {
				// if not, it means that Self has to be greater than Other
				return 1
			}
			break
		}
		lhs := ss[i]
		rhs := so[i]
		if lhs == rhs {
			continue
		} else if lhs < rhs {
			return -1
		}
		return 1
	}
	return 0
}

func (v *Version) LessThan(o *Version) bool {
	return v.Compare(o) < 0
}

func comparePre(v string, other string) int {
	if v == other {
		return 0
	}
	ms := strings.Split(v, ".")
	mo := strings.Split(other, ".")
	ls := len(ms)
	lo := len(mo)
	maxL := lo
	if ls > lo {
		maxL = ls
	}
	for i := 0; i < maxL; i++ {
		sp := ""
		if i < ls {
			sp = ms[i]
		}
		op := ""
		if i < lo {
			op = mo[i]
		}
		compare := comparePart(sp, op)
		if compare != 0 {
			return compare
		}
	}
	return 0
}

func comparePart(v string, other string) int {
	if v == other {
		return 0
	}
	sn := true
	si, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		sn = false
	}
	on := true
	oi, err := strconv.ParseInt(other, 10, 64)
	if err != nil {
		on = false
	}
	if v == "" {
		if on {
			return -1
		}
		return 1
	}
	if other == "" {
		if sn {
			return 1
		}
		return -1
	}
	if sn && !on {
		return -1
	} else if !sn && on {
		return 1
	} else if !sn && !on && v > other {
		return 1
	} else if si > oi {
		return 1
	}
	return -1
}

func allZero(segs []int64) bool {
	for _, s := range segs {
		if s != 0 {
			return false
		}
	}
	return true
}

func GetGitVersion() (*Version, error) {
	return versionCache.Load(nil)
}

func getGitVersion(ctx context.Context) (*Version, error) {
	cmd := command.NewCommand("version")
	result, err := cmd.Run(ctx)
	if err != nil {
		return nil, err
	}
	stdout := result.ReadAsString()
	fields := strings.Fields(stdout)
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid git v output: %s", stdout)
	}
	var v string
	i := strings.Index(fields[2], "windows")
	if i >= 1 {
		v = fields[2][:i-1]
	} else {
		v = fields[2]
	}
	return newVersion(v)
}

func CheckGitVersionAtLeast(v string) error {
	gitVersion, err := GetGitVersion()
	if err != nil {
		return err
	}
	atLeastVersion, err := newVersion(v)
	if err != nil {
		return err
	}
	if gitVersion.Compare(atLeastVersion) < 0 {
		return fmt.Errorf("installed git binary version %s is not at least %s", gitVersion.Original(), v)
	}
	return nil
}
