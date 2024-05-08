package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
)

func Use(engine *gin.Engine) {
	p := ginprometheus.NewPrometheus("golite")
	p.Use(engine)
}
