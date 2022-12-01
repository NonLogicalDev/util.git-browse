package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"EXP/pkg/config"
	"EXP/pkg/gitutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
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

func registerURLFlags(flg *pflag.FlagSet) (service *string, ref *string) {
	service = flg.String(
		"service", "",
		"select service, if other than default",
	)
	ref = flg.String(
		"ref", "",
		"select git ref, if other than default",
	)
	return service, ref
}

func cmdURL() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "url FILEPATH [LineStart] [LineEnd]",
		Short: "output url for a file in remote browser or repo",
	}
	service, ref := registerURLFlags(cmd.PersistentFlags())

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		url, err := getURL(args, service, ref)
		if err != nil {
			return err
		}

		fmt.Println(url)
		return nil
	}
	return cmd
}

func cmdOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open FILEPATH [LineStart] [LineEnd]",
		Short: "open the url for a file in remote browser or repo",
	}
	service, ref := registerURLFlags(cmd.PersistentFlags())

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		url, err := getURL(args, service, ref)
		if err != nil {
			return err
		}
		if len(url) != 0 {
			openInBrowser(url)
		}
		return nil
	}
	return cmd
}

func cmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "print current config",
	}

	flg := cmd.PersistentFlags()
	format := flg.String("format", "gitconfig", "config format (gitconfig | cli | json). Default: gitconfig")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawCfg, err := loadRawConfig()
		if err != nil {
			return err
		}

		cfg, err := config.ParseGitConfig(rawCfg)
		if err != nil {
			return err
		}

		switch *format {
		case "gitconfig":
			fmt.Println(cfg.Dump(config.DumpGitconfig))
		case "cli":
			fmt.Println(cfg.Dump(config.DumpCli))
		case "json":
			fmt.Println(cfg.Dump(config.DumpJson))
		default:
			return fmt.Errorf("unknown config format %v", *format)
		}

		return nil
	}
	return cmd
}

func getURL(args []string, service *string, ref *string) (string, error) {
	args = getArgs(3, args)

	path := args[0]
	lineS, _ := strconv.Atoi(args[1])
	lineE, _ := strconv.Atoi(args[2])

	if len(path) == 0 {
		path = "."
	}

	relPath, isDir, err := gitutils.GetPathInfo(path)
	if err != nil {
		return "", xerrors.Errorf("failed fetching git relative path: %w", err)
	}

	rawCfg, err := loadRawConfig()
	if err != nil {
		return "", xerrors.Errorf("failed loading git configuration: %w", err)
	}

	cfg, err := config.ParseGitConfig(rawCfg)
	if err != nil {
		return "", xerrors.Errorf("failed parsing git configuration: %w", err)
	}

	url, err := cfg.URLFor(config.URLOpt{
		Service: *service,
		Ref:     *ref,
		Path:    relPath,

		IsDir: isDir,

		LineS: lineS,
		LineE: lineE,
	})
	if err != nil {
		return "", xerrors.Errorf("failed rendering url: %w", err)
	}
	return url, nil
}

func loadRawConfig() (string, error) {
	return gitutils.GitExec("config", "--get-regexp", "browse\\..*")
}

func main() {
	log.Default().SetOutput(ioutil.Discard)
	cmd := cmdRoot()
	err := cmd.Execute()
	if err != nil {
		cmd.PrintErrf("TRACE: %+v", err)
		os.Exit(1)

	}
}
