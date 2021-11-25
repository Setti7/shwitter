package form

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type Paginator struct {
	PageState  string `form:"page_state" json:"page_state"`
	PerPage    int    `form:"per_page,default=30" binding:"numeric,gte=1,lte=100" json:"-"`
	NumResults int    `json:"num_results"`
}

func BindPaginatorOrAbort(c *gin.Context) (p *Paginator, err error) {
	err = c.BindQuery(&p)
	return
}

func (p *Paginator) SetResults(iter *gocql.Iter) {
	p.PageState = base64.StdEncoding.EncodeToString(iter.PageState())
	p.NumResults = iter.NumRows()
}

func (p *Paginator) PaginateQuery(q *gocql.Query) *gocql.Query {
	return q.PageState(p.pageStateBytes()).PageSize(p.PerPage)
}

func (p *Paginator) pageStateBytes() []byte {
	data, _ := base64.StdEncoding.DecodeString(p.PageState)
	return data
}
