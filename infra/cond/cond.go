package cond

import (
	"errors"
	"fmt"
)

type Cond struct {
	Field    string      `json:"field"`
	DataType Datatype    `json:"datatype"`
	Op       OPTYPE      `json:"op"`
	Value    interface{} `json:"value"`
}

const PAGESIZEMAX = 1024 * 1024

// 转换
func (c Cond) Build() (sql string, value interface{}, err error) {
	err = nil
	if c.Op == OPEQ {
		return fmt.Sprintf("%s = ?", c.Field), c.Value, err
	} else if c.Op == OPEGT {
		return fmt.Sprintf("%s >= ?", c.Field), c.Value, err
	} else if c.Op == OPGT {
		return fmt.Sprintf("%s > ?", c.Field), c.Value, err
	} else if c.Op == OPLT {
		return fmt.Sprintf("%s < ?", c.Field), c.Value, err
	} else if c.Op == OPLET {
		return fmt.Sprintf("%s <= ?", c.Field), c.Value, err
	} else if c.Op == OPLIKE {
		return fmt.Sprintf("%s like ?", c.Field), "%" + fmt.Sprintf("%s", c.Value) + "%", err
	} else if c.Op == OPIN {
		return fmt.Sprintf("%s in ?", c.Field), c.Value, err
	} else if c.Op == OPBETWEEN {
		if c.DataType == INTARR {
			value = 1
			err = nil
			switch c.Value.(type) {
			case []int:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.Field, c.Value.([]int)[0], c.Value.([]int)[1])
			case []int16:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.Field, c.Value.([]int16)[0], c.Value.([]int16)[1])
			case []int32:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.Field, c.Value.([]int32)[0], c.Value.([]int32)[1])
			case []int64:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.Field, c.Value.([]int64)[0], c.Value.([]int64)[1])
			case []uint:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.Field, c.Value.([]uint)[0], c.Value.([]uint)[1])
			case []float32:
				sql = fmt.Sprintf("%s between %0f and %0f and 1 = ?", c.Field, c.Value.([]float32)[0], c.Value.([]float32)[1])
			case []float64:
				sql = fmt.Sprintf("%s between %0f and %0f and 1 = ?", c.Field, c.Value.([]float64)[0], c.Value.([]float64)[1])
			default:
				sql = ""
				err = errors.New("数据格式不正确")
			}
			return sql, value, err
		} else if c.DataType == STRARR {
			if value, is := c.Value.([]string); is {
				return fmt.Sprintf("%s between %s and %s and 1 = ?", c.Field, value[0], value[1]), 1, err
			} else {
				return "", "", errors.New("数据格式不正确")
			}

		} else {
			if value, is := c.Value.([]string); is {
				return fmt.Sprintf("%s between %s and %s and 1 = ?", c.Field, value[0], value[1]), 1, err
			} else if value, is := c.Value.([]int); is {
				return fmt.Sprintf("%s between %d and %d and 1 = ?", c.Field, value[0], value[1]), 1, err
			} else {
				return " 1 = ? ", 1, err
			}
		}
	} else {
		return " 1 = ? ", 1, err
	}
}

type CondWraper struct {
	Pager Pager  `json:"pager"`
	Order Order  `json:"order"`
	Conds []Cond `json:"conds"`
}

func NewListAllWraper() *CondWraper {
	return &CondWraper{
		Conds: make([]Cond, 0),
		Order: Order{},
		Pager: Pager{Pagefrom: 0, Pagesize: PAGESIZEMAX},
	}
}
func NewCondWrapper() *CondWraper {
	return &CondWraper{
		Conds: make([]Cond, 0),
		Order: Order{},
		Pager: Pager{Pagefrom: 0, Pagesize: 20},
	}
}
func (c *CondWraper) AddCond(cond ...Cond) *CondWraper {
	c.Conds = append(c.Conds, cond...)
	return c
}
func (c *CondWraper) AddOneCond(field string, op OPTYPE, value any) *CondWraper {
	c.Conds = append(c.Conds, Cond{
		Field: field, Op: op, Value: value,
	})
	return c
}
func (c *CondWraper) SetPager(pager Pager) *CondWraper {
	c.Pager = pager
	return c
}

func (c *CondWraper) Pagesize(size int) *CondWraper {
	c.Pager.Pagesize = size
	return c
}
func (c *CondWraper) Pagefrom(from int) *CondWraper {
	c.Pager.Pagefrom = from
	return c
}

func (c *CondWraper) SetOrer(order Order) *CondWraper {
	c.Order = order
	return c
}
