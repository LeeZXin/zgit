package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/common"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"zgit/pkg/git"
	"zgit/pkg/hook"
	"zgit/util"
)

// subHookPreReceive 可用于仓库大小检查提交权限和分支
var subHookPreReceive = &cli.Command{
	Name:        "pre-receive",
	Usage:       "pre-receive Git hook",
	Description: "This command should only be called by Git",
	Action:      runPreReceive,
}

// subHookPostReceive 用于发送通知等
var subHookPostReceive = &cli.Command{
	Name:        "post-receive",
	Usage:       "post-receive Git hook",
	Description: "This command should only be called by Git",
	Action:      runHookPostReceive,
}

var Hook = &cli.Command{
	Name:        "hook",
	Usage:       "This command for zgit hook",
	Description: "zgit",
	Subcommands: []*cli.Command{
		subHookPreReceive,
		subHookPostReceive,
	},
}

func runPreReceive(c *cli.Context) error {
	if isInternal, _ := strconv.ParseBool(os.Getenv(git.EnvIsInternal)); isInternal {
		return nil
	}
	ctx, cancel := initWaitContext(c.Context)
	defer cancel()
	fmt.Println("Welcome to ZGIT")
	// 获取仓库大小限制
	repoId := os.Getenv(git.EnvRepoID)
	repoPath := os.Getenv(git.EnvRepoPath)
	limitSize, err := getRepoLimitSize(repoId)
	if err != nil {
		return err
	}
	repoSize, err := git.GetRepoSize(repoPath)
	if err != nil {
		return err
	}
	if limitSize < repoSize {
		fmt.Printf("checking repo size: %s\n", util.VolumeReadable(repoSize))
		return fmt.Errorf("repo size exceeded limit")
	}
	return scanStdinAndDoHttp(ctx, hook.ApiPreReceiveUrl)
}

// scanStdinAndDoHttp 处理输入并发送http
func scanStdinAndDoHttp(ctx context.Context, httpUrl string) error {
	infoList := make([]hook.RevInfo, 0)
	// the environment is set by serv command
	isWiki, _ := strconv.ParseBool(os.Getenv(git.EnvRepoIsWiki))
	pusherId := os.Getenv(git.EnvPusherID)
	repoId := os.Getenv(git.EnvRepoID)
	prId := os.Getenv(git.EnvPRID)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(string(scanner.Bytes()))
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		refName := git.RefName(fields[2])
		if refName.IsBranch() || refName.IsTag() {
			infoList = append(infoList, hook.RevInfo{
				OldCommitId: fields[0],
				NewCommitId: fields[1],
				RefName:     fields[2],
			})
		}
	}
	client := newHttpClient()
	defer client.CloseIdleConnections()
	partitionList := listutil.Partition(infoList, 30)
	for _, partition := range partitionList {
		reqVO := hook.OptsReqVO{
			RevInfoList: partition,
			IsWiki:      isWiki,
			PusherId:    pusherId,
			RepoId:      repoId,
			PrId:        prId,
		}
		if err := doHttp(ctx, client, reqVO, httpUrl); err != nil {
			return fmt.Errorf("do internal api failed: %v", err)
		}
	}
	return nil
}

func runHookPostReceive(c *cli.Context) error {
	if isInternal, _ := strconv.ParseBool(os.Getenv(git.EnvIsInternal)); isInternal {
		return nil
	}
	ctx, cancel := initWaitContext(c.Context)
	defer cancel()
	return scanStdinAndDoHttp(ctx, hook.ApiPostReceiveUrl)
}

func doHttp(ctx context.Context, client *http.Client, reqVO hook.OptsReqVO, url string) error {
	resp := hook.HttpRespVO{}
	err := httputil.Post(ctx, client,
		fmt.Sprintf("http://localhost:%d/%s", common.DefaultHttpServerPort, url),
		nil,
		reqVO,
		&resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return errors.New(resp.Message)
	}
	return nil
}

func newHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        0, // 禁用连接池
			MaxIdleConnsPerHost: 0, // 禁用连接池
			IdleConnTimeout:     0, // 禁用连接池
		},
		Timeout: 5 * time.Second,
	}
}

func getRepoLimitSize(repoId string) (int64, error) {
	return 1 * util.Gib, nil
}
