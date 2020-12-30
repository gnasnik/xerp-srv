package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/frankffenn/xerp-srv/go-utils/config"
	"github.com/frankffenn/xerp-srv/go-utils/db"
	"github.com/frankffenn/xerp-srv/go-utils/log"
	ursmod "github.com/frankffenn/xerp-srv/services/users/mod"
	"github.com/go-xorm/xorm"
	"golang.org/x/xerrors"
)

var defaultFile = "config.toml"

var (
	App *AppConfig
	eng *xorm.Engine
)

type AppConfig struct {
	Mode          string
	ListenAddress string
	Timeout       int
	JWTScrect     string
	DBUrl         string
}

func InitConfig() error {
	dir, _ := os.Getwd()
	cf := filepath.Join(dir, defaultFile)
	_, err := os.Stat(cf)
	if err != nil {
		return xerrors.Errorf("stat config file (%s): %w", cf, err)
	}
	ff, err := config.FromFile(cf, &AppConfig{})
	if err != nil {
		return xerrors.Errorf("loading config: %w", err)
	}
	App = ff.(*AppConfig)

	logLevel := "info"
	if strings.ToUpper(App.Mode) == "DEBUG" {
		logLevel = "debug"
	}

	ld := filepath.Join(dir, "log")
	_, err = os.Stat(ld)
	if os.IsNotExist(err) {
		err = os.MkdirAll(ld, os.ModePerm)
	}
	if err != nil {
		return xerrors.Errorf("stat log dir err: %w", err)
	}

	logger, _ := log.NewLogger(filepath.Join(filepath.Base(ld), os.Args[0]+".log"), logLevel)
	log.SetDefault(logger)

	return nil
}

func initDB() error {
	var err error
	eng, err = db.OpenDB("mysql", App.DBUrl)
	if err != nil {
		return err
	}
	eng.Sync2(
		new(ursmod.User),
	)
}

// Session ...
func Session() *xorm.Session {
	return eng.NewSession()
}
