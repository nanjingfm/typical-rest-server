package buildtool

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/buildtool/releaser"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const (
	readmeFile = "README.md"
)

type buildtool struct {
	*typictx.Context
}

func (t buildtool) commands() (cmds []*typictx.Command) {
	cmds = []*typictx.Command{
		{
			Name:       "build",
			ShortName:  "b",
			Usage:      "Build the binary",
			ActionFunc: t.buildBinary,
		},
		{
			Name:       "clean",
			ShortName:  "c",
			Usage:      "Clean the project from generated file during build time",
			ActionFunc: t.cleanProject,
		},
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the binary",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-build",
					Usage: "Run the binary without build",
				},
			},
			ActionFunc: t.runBinary,
		},
		{
			Name:       "test",
			ShortName:  "t",
			Usage:      "Run the testing",
			ActionFunc: t.runTesting,
		},
		{
			Name:  "release",
			Usage: "Release the distribution",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-test",
					Usage: "Release without run automated test",
				},
				cli.BoolFlag{
					Name:  "no-github",
					Usage: "Release without create github release",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Release by passed all validation",
				},
				cli.BoolFlag{
					Name:  "alpha",
					Usage: "Release for alpha version",
				},
			},
			ActionFunc: t.releaseDistribution,
		},
		{
			Name:  "mock",
			Usage: "Generate mock class",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-delete",
					Usage: "Generate mock class with delete previous generation",
				},
			},
			ActionFunc: t.generateMock,
		},
		{
			Name:       "readme",
			Usage:      "Generate readme document",
			ActionFunc: t.generateReadme,
		},
		{
			Name:       "docker",
			Usage:      "Docker utility",
			BeforeFunc: typienv.LoadEnvFile,
			SubCommands: []*typictx.Command{
				{
					Name:       "compose",
					Usage:      "Generate docker-compose.yaml",
					ActionFunc: t.dockerCompose,
				},
				{
					Name:  "up",
					Usage: "Create and start containers",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "no-compose",
							Usage: "Create and start containers without generate docker-compose.yaml",
						},
					},
					ActionFunc: t.dockerUp,
				},
				{
					Name:       "down",
					Usage:      "Stop and remove containers, networks, images, and volumes",
					ActionFunc: t.dockerDown,
				},
			},
		},
	}
	cmds = append(cmds, t.CommandLines()...)
	return
}

func (t buildtool) buildBinary(ctx *typictx.ActionContext) error {
	log.Info("Build application binary")
	return bash.GoBuild(typienv.App.BinPath, typienv.App.SrcPath)
}

func (t buildtool) cleanProject(ctx *typictx.ActionContext) (err error) {
	log.Info("Start clean the project")
	log.Infof("  Remove %s", typienv.Bin)
	if err = os.RemoveAll(typienv.Bin); err != nil {
		return
	}
	log.Infof("  Remove %s", typienv.Metadata)
	if err = os.RemoveAll(typienv.Metadata); err != nil {
		return
	}
	log.Info("  Remove .env")
	os.Remove(".env")
	return filepath.Walk(typienv.Dependency.SrcPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			log.Infof("  Remove %s", path)
			return os.Remove(path)
		}
		return nil
	})
}

func (t buildtool) runBinary(ctx *typictx.ActionContext) (err error) {
	if !ctx.Cli.Bool("no-build") {
		if err = t.buildBinary(ctx); err != nil {
			return
		}
	}
	log.Info("Run application binary")
	return bash.Run(typienv.App.BinPath, []string(ctx.Cli.Args())...)
}

func (t buildtool) runTesting(ctx *typictx.ActionContext) error {
	log.Info("Run testings")
	return bash.GoTest(ctx.TestTargets)
}

func (t buildtool) generateMock(ctx *typictx.ActionContext) (err error) {
	log.Info("Generate mocks")
	if err = bash.GoGet("github.com/golang/mock/mockgen"); err != nil {
		return
	}
	mockPkg := typienv.Mock
	if !ctx.Cli.Bool("no-delete") {
		log.Infof("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}
	for _, mockTarget := range ctx.MockTargets {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]
		err = bash.RunGoBin("mockgen",
			"-source", mockTarget,
			"-destination", dest,
			"-package", mockPkg)
	}
	return
}

func (t buildtool) releaseDistribution(action *typictx.ActionContext) (err error) {
	log.Info("Release distribution")
	var binaries, changeLogs []string
	if !action.Cli.Bool("no-test") {
		if err = t.runTesting(action); err != nil {
			return
		}
	}
	force := action.Cli.Bool("force")
	alpha := action.Cli.Bool("alpha")
	if binaries, changeLogs, err = releaser.ReleaseDistribution(action.Release, force, alpha); err != nil {
		return
	}
	if !action.Cli.Bool("no-github") {
		releaser.GithubRelease(binaries, changeLogs, action.Release, alpha)
	}
	return
}

func (t buildtool) dockerCompose(ctx *typictx.ActionContext) (err error) {
	log.Info("Generate docker-compose.yml")
	dockerCompose := ctx.DockerCompose()
	d1, _ := yaml.Marshal(dockerCompose)
	return ioutil.WriteFile("docker-compose.yml", d1, 0644)
}

func (t buildtool) dockerUp(ctx *typictx.ActionContext) (err error) {
	cmd := exec.Command("docker-compose", "up", "--remove-orphans", "-d")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (t buildtool) dockerDown(ctx *typictx.ActionContext) (err error) {
	cmd := exec.Command("docker-compose", "down")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (t buildtool) generateReadme(a *typictx.ActionContext) (err error) {
	var file *os.File
	log.Infof("Generate Readme: %s", readmeFile)
	if file, err = os.Create(readmeFile); err != nil {
		return
	}
	defer file.Close()
	Readme{
		Title:       a.Name,
		Description: a.Description,
		Context:     a.Context,
	}.Markdown(file)
	return
}