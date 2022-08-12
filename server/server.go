package server

import (
	"github.com/brbarme-shop/brbarmex-rating/server/route"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	route.LoadRoute(r)

	r.Run("localhost:3001")
}
