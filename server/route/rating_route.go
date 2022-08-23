package route

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/brbarme-shop/brbarmex-rating/config"
	"github.com/brbarme-shop/brbarmex-rating/postgresql"
	"github.com/brbarme-shop/brbarmex-rating/rating"
	"github.com/gin-gonic/gin"
)

var (
	cfg              = config.NewConfiguration()
	db               = postgresql.NewSqlDB(cfg)
	ratingRepository = postgresql.NewRatingRepository(db)
)

func postRating(c *gin.Context) {

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
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
}
