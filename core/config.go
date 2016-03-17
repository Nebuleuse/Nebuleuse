package core

import (
	"github.com/robfig/config"
	"strconv"
)

var defaultOptions = map[string]string{"serverAddress": "127.0.0.1",
	"serverPort":               "8080",
	"dbType":                   "mysql",
	"dbAddress":                "127.0.0.1:3306",
	"dbUser":                   "",
	"dbPass":                   "",
	"dbBase":                   "",
	"MaxSessionsChannelBuffer": "10",
	"LongpollingTimeout":       "10",
	"DashboardLocation":        "./admin/"}

func (c *configMgr) InitConfig() error {
	var err error
	c.Cfg = make(map[string]string)
	c.SysCfg = make(map[string]string)

	configFile, err := config.ReadDefault(".config")
	if err != nil {
		Error.Fatal("Could not open config file. Please run 'nebuleuse install' to generate one")
		return err
	}
	configs, err := configFile.Options("default")
	if err != nil {
		Error.Fatal("Could not read values from config file")
		return err
	}

	for _, value := range configs {
		c.SysCfg[value], err = configFile.String("default", value)
	}
	return nil
}
func (c *configMgr) GetSysConfig(name string) string {
	c.configLock.RLock()
	res := c.SysCfg[name]
	c.configLock.RUnlock()
	return res
}
func (c *configMgr) GetSysConfigInt(name string) int {
	c.configLock.RLock()
	res, _ := strconv.Atoi(c.SysCfg[name])
	c.configLock.RUnlock()
	return res
}
func (c *configMgr) GetConfig(name string) string {
	c.configLock.RLock()
	res := c.Cfg[name]
	c.configLock.RUnlock()
	return res
}
func (c *configMgr) GetConfigInt(name string) int {
	c.configLock.RLock()
	res, _ := strconv.Atoi(c.Cfg[name])
	c.configLock.RUnlock()
	return res
}
func (c *configMgr) GetConfigFloat(name string) float64 {
	c.configLock.RLock()
	res, _ := strconv.ParseFloat(c.Cfg[name], 64)
	c.configLock.RUnlock()
	return res
}

func (c *configMgr) LoadConfig() {
	var (
		name  string
		value string
	)

	rows, err := Db.Query("select name, value from neb_config")
	if err != nil {
		Error.Fatal(err)
	}

	defer rows.Close()
	c.configLock.Lock()
	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			Error.Fatal(err)
		}
		c.Cfg[name] = value
	}
	c.configLock.Unlock()
	err = rows.Err()
	if err != nil {
		Error.Fatal(err)
	} else {
		Info.Println("Successfully read configuration")
	}
}

func (c *configMgr) SetConfig(name, value string) error {
	c.configLock.Lock()
	_, err := Db.Query("UPDATE neb_config SET value=? WHERE name=?", value, name)
	if err != nil {
		Error.Println("Failed to update config : ", value, name)
		return err
	}
	c.Cfg[name] = value
	c.configLock.Unlock()
	return nil
}
