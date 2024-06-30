package cond

type OPTYPE string

const (
	OPEQ      OPTYPE = "eq"
	OPGT      OPTYPE = "gt"
	OPEGT     OPTYPE = "egt"
	OPLT      OPTYPE = "lt"
	OPLET     OPTYPE = "let"
	OPLIKE    OPTYPE = "like"
	OPBETWEEN OPTYPE = "between"
	OPIN      OPTYPE = "in"
)
