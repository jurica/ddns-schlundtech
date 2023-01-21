package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"

	"github.com/jurica/ddns-schlundtech/schlundtech"
)

var router = gin.Default()
var config Config

type Config struct {
	User     string
	Password string
	Context  string
}

func main() {
	rawConfig, err := os.ReadFile("ddns-schlundtech.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(rawConfig, &config)
	if err != nil {
		panic(err)
	}
	if config.Context == "" || config.Password == "" || config.User == "" {
		panic("please provide values for all required authentication keys: user, password, context")
	}

	router.GET("", serve)

	router.Run()
}

func serve(c *gin.Context) {
	// TODO get data from env and req
	domain := c.Query("domain")
	ip := c.Query("ip")

	if ip == "" || domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please provide both query parameters 'domain' and 'ip'",
		})

		return
	}

	err := schlundtech.UpdateDdnsRecord(config.User, config.Password, config.Context, "", domain, ip)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}
