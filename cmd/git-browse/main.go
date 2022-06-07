package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"

	"EXP/pkg/config"
	"EXP/pkg/gitutils"
)

func getArgs(n int, args []string) []string {
	out := make([]string, n)
	for i := 0; i < len(args) && i < n; i++ {
		out[i] = args[i]
	}
	return out
}

func cmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use: "git-browse",
	}
	cmd.AddCommand(cmdURL())
	cmd.AddCommand(cmdOpen())
	cmd.AddCommand(cmdConfig())
	return cmd
}

func cmdURL() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "url [FILEPATH] [LineStart] [LineEnd]",
		Short: "output url for a file in remote browser or repo",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		args = getArgs(3, args)

		path := args[0]
		lineS, _ := strconv.Atoi(args[1])
		lineE, _ := strconv.Atoi(args[2])

		if len(path) == 0 {
			path = "."
		}
		relPath, _ := gitutils.GetRelativePath(path)

		isDir := false
		pathStat, err := os.Stat(path)
		if err == nil {
			isDir = pathStat.IsDir()
		}

		cfg, _ := config.ParseConfig()
		url, err := cfg.GetURL("", relPath, isDir, lineS, lineE)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println(url)
	}
	return cmd
}

func cmdOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "open the url for a file in remote browser or repo",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		args = getArgs(3, args)

		path := args[0]
		lineS, _ := strconv.Atoi(args[1])
		lineE, _ := strconv.Atoi(args[2])

		if len(path) == 0 {
			path = "."
		}
		relPath, _ := gitutils.GetRelativePath(path)

		isDir := false
		pathStat, err := os.Stat(path)
		if err == nil {
			isDir = pathStat.IsDir()
		}

		cfg, _ := config.ParseConfig()
		url, err := cfg.GetURL("", relPath, isDir, lineS, lineE)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		if len(url) != 0 {
			openInBrowser(url)
		}
	}
	return cmd
}

func cmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "print current config",
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		cfg, _ := config.ParseConfig()
		litter.Dump(cfg)
	}
	return cmd
}

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openInBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		fmt.Println("openning...", url)
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		fmt.Println("openning...", url)
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		fmt.Println("openning...", url)
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Default().SetOutput(ioutil.Discard)
	_ = cmdRoot().Execute()
}
