package core

import (
	"github.com/robfig/config"
)

func initConfig() error {
	var err error
	Cfg = make(map[string]string)
	SysCfg = make(map[string]string)

	cfgNames := [...]string{"serverAddress", "serverPort", "dbType", "dbAddress", "dbUser", "dbPass", "dbBase", "gitPath"}
	c, err := config.ReadDefault(".config")
	for _, v := range cfgNames {
		SysCfg[v], err = c.String("default", v)
		if err != nil {
			Error.Fatal("Could not read config: " + v)
			return err
		}
	}
	return nil
}

func loadConfig() {
	var (
		name  string
		value string
	)

	rows, err := Db.Query("select name, value from neb_config")
	if err != nil {
		Error.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			Error.Fatal(err)
		}
		Cfg[name] = value
	}

	err = rows.Err()
	if err != nil {
		Error.Fatal(err)
	} else {
		Info.Println("Successfully read configuration")
	}
}

func (c *ConfigMgr) SetConfig(name, value string) error {
	_, err := Db.Query("UPDATE neb_config SET value=? WHERE name=?", value, name)
	if err != nil {
		return err
	}
	(*c)[name] = value
	return nil
}
