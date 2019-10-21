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

	// if getAppEnv() == "test" {

	// } else {
	driver := config.Conf.Get("database.driver").(string)
	configTree := config.Conf.Get(driver).(*toml.Tree)

	connect := configTree.Get("connect").(string)

	DB, err := gorm.Open(driver, connect)

	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}

	return DB
	// }
}
