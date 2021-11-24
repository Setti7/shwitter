package form

import (
	"github.com/gin-gonic/gin"
)

// TODO: move this to a module?
type Paginator struct {
	Ref        string `form:"ref"`
	NumResults int    `form:"num_results,default=30" binding:"numeric,gte=1,lte=100"`
}

func BindPaginatorOrAbort(c *gin.Context) (p *Paginator, err error) {
	err = c.BindQuery(&p)
	return
}
