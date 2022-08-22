package route

import (
	"encoding/json"
	"fmt"
	"io"
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

	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
	}

	defer c.Request.Body.Close()

	var ratingInput *rating.PutRatingInput
	err = json.Unmarshal(b, &ratingInput)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
	}

	err = rating.PutRating(c.Request.Context(), ratingInput, ratingRepository)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
}
