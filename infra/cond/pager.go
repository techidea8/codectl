package cond

type Pager struct {
	Pagefrom int `json:"pagefrom"`
	Pagesize int `json:"pagesize"`
}

func (p Pager) Limit() int {
	return p.Pagesize
}
func (p Pager) Offset() int {
	if p.Pagefrom < 1 {
		return 0
	} else {
		return (p.Pagefrom - 1) * p.Pagesize
	}
}
