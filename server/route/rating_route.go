package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Post
func addRating(c *gin.Context) {

	var anything interface{}

	if err := c.BindJSON(&anything); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, anything)
}
