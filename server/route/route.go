package route

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type routeMap struct {
	methodHTTP string
	uri        string
	action     func(gin *gin.Context)
}

var routeMaps = []routeMap{
	{
		methodHTTP: http.MethodPost,
		uri:        "/rating",
		action:     addRating,
	},
}

func LoadRoute(r *gin.Engine) {

	for _, rm := range routeMaps {

		switch rm.methodHTTP {
		case http.MethodPost:
			r.POST(rm.uri, rm.action)

		default:
			log.Println("pau")
		}

	}
}
