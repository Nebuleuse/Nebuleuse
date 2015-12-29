package core

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func InitMock(t *testing.T, cfg, sysCfg ConfigMgr, sessions, messaging, git bool) sqlmock.Sqlmock {
	SysCfg = sysCfg
	cfg = cfg

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	Db = db

	if sessions {
		initSessions()
	}
	if messaging {
		initMessaging()
	}
	if git {
		if Cfg["updateSystem"] == "GitPatch" || Cfg["updateSystem"] == "FullGit" {
			initGit()
		}
	}
	return mock
}
