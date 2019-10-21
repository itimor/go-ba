package database

import (
	"fmt"

	"../config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pelletier/go-toml"
)

var (
	DB = InitDB()
)

/**
 * Set up a database connection
 * @param diver string
 */

func InitDB() *gorm.DB {

	getAppEnv := "test"
	// if getAppEnv() == "test" {
	// } else {}

	configTree := config.Conf.Get(getAppEnv + ".database").(*toml.Tree)
	driver := configTree.Get("driver").(string)
	connect := configTree.Get("connect").(string)

	DB, err := gorm.Open(driver, connect)

	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}

	return DB
	// }
}
