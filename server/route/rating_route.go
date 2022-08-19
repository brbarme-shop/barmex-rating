package route

import (
	"net/http"

	"github.com/brbarme-shop/brbarmex-rating/postgresql"
	"github.com/brbarme-shop/brbarmex-rating/rating"
	"github.com/gin-gonic/gin"
)

var (
	db               = postgresql.NewSqlDB("host=localhost port=5432 user=rating_user password=rating_pwd dbname=rating_db sslmode=disable")
	ratingRepository = postgresql.NewRatingRepository(db)
)

// @Post
func addRating(c *gin.Context) {

	var ratingInput *rating.RatingInput

	err := c.BindJSON(&ratingInput)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	err = rating.PutRating(c.Request.Context(), ratingInput, ratingRepository)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
}
