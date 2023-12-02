package cmd

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf/starter"
	"github.com/LeeZXin/zsf/zsf"
	"github.com/urfave/cli/v2"
	"zgit/git"
	"zgit/ssh"
	"zgit/util"
)

var Web = &cli.Command{
	Name:        "web",
	Usage:       "This command starts zgit server",
	Description: "zgit",
	Action:      runWeb,
}

func runWeb(c *cli.Context) error {
	options := make([]zsf.Option, 0)
	args := c.Args().Slice()
	for _, arg := range args {
		if arg == "debug" {
			options = append(options, zsf.WithRunMode("debug"))
		}
	}
	ssh.InitSsh()
	git.InitGit()
	fmt.Println(git.Merge(context.Background(), util.JoinRepoPath("lizexin", "ppp-test"), "dev", "master", git.MergeRepoOpts{
		Message:  "fuck uuuu",
		SigKeyId: "",
	}))
	// 6412c3c284ed5be776e3c8f800c608a04a003f61 142ad803fee8ad124c554f141994849163f7e8ae "clue-service/src/main/java/cn/wesure/clue/strategy/clue/handlerChain/InterruptClueHandler.java"
	//content, err := git.PreparePullRequest(context.Background(), "D:\\Projects\\wecare-clue-runner", "dev_4.45.0_1206", "master")
	//if err != nil {
	//	panic(err)
	//} else {
	//	marshal, _ := json.Marshal(content)
	//	util.WriteFile("ggg.txt", marshal)
	//}
	starter.Run(options...)
	return nil
}
