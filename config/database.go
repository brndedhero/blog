package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

var DB *gorm.DB

func dbDsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		3306,
		os.Getenv("DB_NAME"),
	)
}

func ConnectDb() *gorm.DB {
	con, err := gorm.Open(mysql.Open(dbDsn()), &gorm.Config{})
	if err != nil {
		Log.WithFields(logrus.Fields{
			"app":  "mysql",
			"func": "connectDb",
		}).Fatal(err)
	}
	con.Use(prometheus.New(prometheus.Config{
		DBName:          os.Getenv("DB_NAME"),
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		},
	}))
	return con
}
