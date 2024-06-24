package gin

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryArrayInt64(c *gin.Context, key string) (idN64s []int64, err error) {
	ids := c.QueryArray(key)
	for _, id := range ids {
		var idN64 int64
		idN64, err = strconv.ParseInt(id, 10, 16)
		if err != nil {
			return
		}
		idN64s = append(idN64s, idN64)
	}
	return
}

func QueryInt64(c *gin.Context, key string) (value int64, err error) {
	value, err = strconv.ParseInt(c.Query(key), 10, 16)
	if err != nil {
		return
	}
	return
}
