package route

import (
	"net/http"

	"github.com/brbarme-shop/brbarmex-rating/postgresql"
	"github.com/brbarme-shop/brbarmex-rating/rating"
	"github.com/gin-gonic/gin"
)

var (
	db = postgresql.NewRatingRepository("host=localhost port=5432 user=rating_user password=rating_pwd dbname=rating_db sslmode=disable")
)

// @Post
func addRating(c *gin.Context) {

	var ratingInput *rating.AverageInput

	err := c.BindJSON(&ratingInput)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	err = rating.PutRatingAverage(c.Request.Context(), ratingInput, db)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, ratingInput)
}
