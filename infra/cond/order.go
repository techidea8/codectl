package cond

type Order struct {
	Field  string `json:"field"`
	Method string `json:"method"`
}

func (o Order) Desc() string {
	return o.Field + " desc "
}
func (o Order) Asc() string {
	return o.Field + " asc "
}
func (o Order) Build() (string, error) {
	if o.Method == "asc" {
		return o.Asc(), nil
	} else if o.Method == "desc" {
		return o.Desc(), nil
	} else {
		return "", nil
	}
}
