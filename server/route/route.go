package route

import (
	"fmt"
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
	{
		methodHTTP: http.MethodGet,
		uri:        "/health",
		action:     healthCheck,
	},
}

func LoadRoute(r *gin.Engine) {

	for _, rm := range routeMaps {

		switch rm.methodHTTP {
		case http.MethodPost:
			r.POST(rm.uri, rm.action)
		case http.MethodGet:
			r.GET(rm.uri, rm.action)
		default:
			panic(fmt.Sprintln("method http not implemented or is invalid ", rm.methodHTTP))
		}

	}
}
