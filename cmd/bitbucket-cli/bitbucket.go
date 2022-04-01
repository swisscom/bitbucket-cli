package main

import (
	"github.com/adrg/xdg"
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	cli "github.com/swisscom/bitbucket-cli/internal"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

type Args struct {
	Debug       *bool           `arg:"-D,--debug"`
	Username    string          `arg:"-u,--username,env:BITBUCKET_USERNAME"`
	Password    string          `arg:"-p,--password,env:BITBUCKET_PASSWORD"`
	AccessToken string          `arg:"-t,--access-token,env:BITBUCKET_ACCESS_TOKEN" help:"A Personal Access Token"`
	Url         string          `arg:"-u,--url,env:BITBUCKET_URL" help:"URL to the REST API of Bitbucket, e.g: https://git.example.com/rest"`
	Config      string          `arg:"-c,--config"`
	Project     *cli.ProjectCmd `arg:"subcommand:project"`
	Repo        *cli.RepoCmd    `arg:"subcommand:repo"`
	Pr          *cli.PrCmd      `arg:"subcommand:pr"`
}

var args Args

func main() {
	p := arg.MustParse(&args)
	logger := logrus.New()

	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	loadConfigIntoArgs(logger)

	if args.Url == "" {
		logger.Fatalf("please specifiy a Bitbucket URL")
	}

	if args.Username == "" {
		logger.Fatalf("An username is required")
	}

	var auth cli.Authenticator
	if args.Password != "" {
		// Basic Auth
		auth = cli.BasicAuth{Username: args.Username, Password: args.Password}
	} else if args.AccessToken != "" {
		auth = cli.AccessToken{Username: args.Username, AccessToken: args.AccessToken}
	} else {
		logger.Fatalf("either a password or an access token must be provided")
	}

	c, err := cli.NewCLI(auth, args.Url)
	if err != nil {
		logger.Fatalf("unable to create CLI: %v", err)
	}
	c.SetLogger(logger)

	if args.Project != nil {
		c.RunProjectCmd(args.Project)
		return
	}

	if args.Repo != nil {
		c.RunRepoCmd(args.Repo)
		return
	}

	if args.Pr != nil {
		c.RunPRCmd(args.Pr)
		return
	}

	p.Fail("Command must be specified")

}

func loadConfigIntoArgs(logger *logrus.Logger) {
	var cfg *cli.Config
	var err error

	if args.Config != "" {
		// Load config if specified
		cfg, err = loadConfig(logger, args.Config)
		if err != nil {
			logger.Fatalf("unable to load config: %v", err)
		}
	} else {
		stdCfgFilePath := path.Join(xdg.ConfigHome, "bitbucket-cli", "config.yml")
		logger.Debugf("loading config from %v", stdCfgFilePath)
		_, err = os.Stat(stdCfgFilePath)
		if err == nil {
			// File exists, let's load the config
			cfg, err = loadConfig(logger, stdCfgFilePath)
			if err != nil {
				logger.Fatalf("cannot parse config (%s): %v", stdCfgFilePath, err)
			}
		}
	}

	if cfg != nil {
		args.Username = setIfEmpty(args.Username, cfg.Username)
		args.Password = setIfEmpty(args.Password, cfg.Password)
		args.AccessToken = setIfEmpty(args.AccessToken, cfg.AccessToken)
		args.Url = setIfEmpty(args.Url, cfg.Url)
	}
}

func setIfEmpty(argValue string, cfgValue string) string {
	if argValue != "" {
		return argValue
	}

	return cfgValue
}

func loadConfig(logger *logrus.Logger, configPath string) (*cli.Config, error) {
	if configPath == "" {
		return nil, nil
	}

	f, err := os.Open(configPath)
	if err != nil {
		logger.Fatalf("unable to open config: %v", err)
	}

	var config cli.Config
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&config)
	if err != nil {
		logger.Fatalf("unable to parse config: %v", err)
	}

	return &config, nil
}
