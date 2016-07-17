package httpapi

import (
	"champ/plex/model"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Navigation struct {
}

func (a *Navigation) Handle(c *gin.Context) {
	log.Info(c.Request.URL.Query())
	c.XML(200, model.Response{Code: "200", Status: "OK"})
}
