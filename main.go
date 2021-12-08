package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/brndedhero/blog/config"
	"github.com/brndedhero/blog/models"
	"github.com/brndedhero/blog/router"
	"github.com/sirupsen/logrus"
)

func main() {
	config.DB = config.ConnectDb()
	config.DB.AutoMigrate(&models.BlogPost{}, &models.Tag{})
	config.Redis = config.ConnectRedis()
	config.Log = config.SetupLogger()

	http.Handle("/", router.SetupRouter())
	httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))
	message := fmt.Sprintf("Listening for requests at http://%s:%d", os.Getenv("HTTP_URL"), httpPort)
	config.Log.WithFields(logrus.Fields{
		"app":  "blog",
		"func": "main",
	}).Info(message)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil))
}
