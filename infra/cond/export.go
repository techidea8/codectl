package cond

type Meta struct {
	Prop  string `json:"prop"`
	Label string `json:"label"`
}
type Export struct {
	Meta []Meta      `json:"meta"`
	Cond *CondWraper `json:"cond"`
}

func NewExport() *Export {
	r := &Export{
		Meta: make([]Meta, 0),
		Cond: NewListAllWraper(),
	}
	return r
}
