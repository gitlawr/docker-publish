package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/rancher/pipeline-docker-publish/docker"
)

var build = "0" // build number set at compile-time

func main() {

	app := cli.NewApp()
	app.Name = "docker publish"
	app.Usage = "docker publish"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "push-remote",
			Usage:  "push to remote docker registry",
			EnvVar: "PUSH_REMOTE",
		},
		cli.StringFlag{
			Name:   "remote.url",
			Usage:  "git remote url",
			EnvVar: "CICD_GIT_URL",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "CICD_GIT_COMMIT",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "daemon.mirror",
			Usage:  "docker daemon registry mirror",
			EnvVar: "MIRROR",
		},
		cli.StringFlag{
			Name:   "daemon.storage-driver",
			Usage:  "docker daemon storage driver",
			EnvVar: "STORAGE_DRIVER",
		},
		cli.StringFlag{
			Name:   "daemon.storage-path",
			Usage:  "docker daemon storage path",
			Value:  "/var/lib/docker",
			EnvVar: "STORAGE_PATH",
		},
		cli.StringFlag{
			Name:   "daemon.bip",
			Usage:  "docker daemon bride ip address",
			EnvVar: "BIP",
		},
		cli.StringFlag{
			Name:   "daemon.mtu",
			Usage:  "docker daemon custom mtu setting",
			EnvVar: "MTU",
		},
		cli.StringSliceFlag{
			Name:   "daemon.dns",
			Usage:  "docker daemon dns server",
			EnvVar: "CUSTOM_DNS",
		},
		cli.StringSliceFlag{
			Name:   "daemon.dns-search",
			Usage:  "docker daemon dns search domains",
			EnvVar: "CUSTOM_DNS_SEARCH",
		},
		cli.BoolFlag{
			Name:   "daemon.insecure",
			Usage:  "docker daemon allows insecure registries",
			EnvVar: "INSECURE",
		},
		cli.BoolFlag{
			Name:   "daemon.ipv6",
			Usage:  "docker daemon IPv6 networking",
			EnvVar: "IPV6",
		},
		cli.BoolFlag{
			Name:   "daemon.experimental",
			Usage:  "docker daemon Experimental mode",
			EnvVar: "EXPERIMENTAL",
		},
		cli.BoolFlag{
			Name:   "daemon.debug",
			Usage:  "docker daemon executes in debug mode",
			EnvVar: "DEBUG,DOCKER_LAUNCH_DEBUG",
		},
		cli.BoolFlag{
			Name:   "daemon.off",
			Usage:  "don't start the docker daemon",
			EnvVar: "DAEMON_OFF",
		},
		cli.StringFlag{
			Name:   "dockerfile",
			Usage:  "build dockerfile",
			Value:  "Dockerfile",
			EnvVar: "DOCKERFILE",
		},
		cli.StringFlag{
			Name:   "context",
			Usage:  "build context",
			Value:  ".",
			EnvVar: "CONTEXT",
		},
		cli.StringSliceFlag{
			Name:   "tags",
			Usage:  "build tags",
			Value:  &cli.StringSlice{"latest"},
			EnvVar: "TAG,TAGS",
		},
		cli.BoolFlag{
			Name:   "tags.auto",
			Usage:  "default build tags",
			EnvVar: "DEFAULT_TAGS,AUTO_TAG",
		},
		cli.StringFlag{
			Name:   "tags.suffix",
			Usage:  "default build tags with suffix",
			EnvVar: "DEFAULT_SUFFIX,AUTO_TAG_SUFFIX",
		},
		cli.StringSliceFlag{
			Name:   "args",
			Usage:  "build args",
			EnvVar: "BUILD_ARGS",
		},
		cli.StringSliceFlag{
			Name:   "args-from-env",
			Usage:  "build args",
			EnvVar: "BUILD_ARGS_FROM_ENV",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "build target",
			EnvVar: "TARGET",
		},
		cli.BoolFlag{
			Name:   "squash",
			Usage:  "squash the layers at build time",
			EnvVar: "SQUASH",
		},
		cli.BoolTFlag{
			Name:   "pull-image",
			Usage:  "force pull base image at build time",
			EnvVar: "PULL_IMAGE",
		},
		cli.BoolFlag{
			Name:   "compress",
			Usage:  "compress the build context using gzip",
			EnvVar: "COMPRESS",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "docker repository",
			EnvVar: "REPO",
		},
		cli.StringSliceFlag{
			Name:   "label-schema",
			Usage:  "label-schema labels",
			EnvVar: "LABEL_SCHEMA",
		},
		cli.StringFlag{
			Name:   "docker.registry",
			Usage:  "docker registry",
			Value:  "https://index.docker.io/v1/",
			EnvVar: "REGISTRY,DOCKER_REGISTRY",
		},
		cli.StringFlag{
			Name:   "docker.username",
			Usage:  "docker username",
			EnvVar: "USERNAME,DOCKER_USERNAME",
		},
		cli.StringFlag{
			Name:   "docker.password",
			Usage:  "docker password",
			EnvVar: "PASSWORD,DOCKER_PASSWORD",
		},
		cli.StringFlag{
			Name:   "docker.email",
			Usage:  "docker email",
			EnvVar: "EMAIL,DOCKER_EMAIL",
		},
		cli.BoolTFlag{
			Name:   "docker.purge",
			Usage:  "docker should cleanup images",
			EnvVar: "PURGE",
		},
		cli.StringFlag{
			Name:   "repo.branch",
			Usage:  "repository branch",
			EnvVar: "CICD_GIT_BRANCH",
		},
		cli.BoolFlag{
			Name:   "no-cache",
			Usage:  "do not use cached intermediate containers",
			EnvVar: "NO_CACHE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := docker.Plugin{
		PushRemote: c.Bool("push-remote"),
		Cleanup:    c.BoolT("docker.purge"),
		Login: docker.Login{
			Registry: c.String("docker.registry"),
			Username: c.String("docker.username"),
			Password: c.String("docker.password"),
			Email:    c.String("docker.email"),
		},
		Build: docker.Build{
			Remote:      c.String("remote.url"),
			Name:        c.String("commit.sha"),
			Dockerfile:  c.String("dockerfile"),
			Context:     c.String("context"),
			Tags:        c.StringSlice("tags"),
			Args:        c.StringSlice("args"),
			ArgsEnv:     c.StringSlice("args-from-env"),
			Target:      c.String("target"),
			Squash:      c.Bool("squash"),
			Pull:        c.BoolT("pull-image"),
			Compress:    c.Bool("compress"),
			Registry:    c.String("docker.registry"),
			Repo:        c.String("repo"),
			LabelSchema: c.StringSlice("label-schema"),
			NoCache:     c.Bool("no-cache"),
		},
		Daemon: docker.Daemon{
			Registry:      c.String("docker.registry"),
			Mirror:        c.String("daemon.mirror"),
			StorageDriver: c.String("daemon.storage-driver"),
			StoragePath:   c.String("daemon.storage-path"),
			Insecure:      c.Bool("daemon.insecure"),
			Disabled:      c.Bool("daemon.off"),
			IPv6:          c.Bool("daemon.ipv6"),
			Debug:         c.Bool("daemon.debug"),
			Bip:           c.String("daemon.bip"),
			DNS:           c.StringSlice("daemon.dns"),
			DNSSearch:     c.StringSlice("daemon.dns-search"),
			MTU:           c.String("daemon.mtu"),
			Experimental:  c.Bool("daemon.experimental"),
		},
	}

	return plugin.Exec()
}
