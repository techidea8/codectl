package xlskit

import (
	"bytes"
	"fmt"

	"github.com/techidea8/codectl/infra/slicekit"
	"github.com/xuri/excelize/v2"
)

type xlskitctrl struct {
	sheetname string
	file      *excelize.File
}

func NewXlsCtrl(sheetname string) *xlskitctrl {
	return &xlskitctrl{
		sheetname: sheetname,
		file:      excelize.NewFile(),
	}
}

// meta: 是元数据,{"字段1","字段2","字符3"}
func (s *xlskitctrl) Render(meta []Meta, dataInput []map[string]any) (buf *bytes.Buffer, err error) {
	index, err := s.file.NewSheet(s.sheetname)
	if err != nil {
		return
	}
	s.file.SetActiveSheet(index)
	titleMap := make(map[string]string, 0)
	// 形成字段名称和A,B,C,D之间的映射关系
	fieldArr := []string{}
	for _, value := range meta {
		titleMap[value.Field] = value.Title
		fieldArr = append(fieldArr, value.Field)
	}
	// 添加表头
	rowIndex := 1
	for index, field := range fieldArr {
		cellName := fmt.Sprintf("%s%d", charfileds[index], rowIndex)
		title := titleMap[field]
		s.file.SetCellValue(s.sheetname, cellName, title)
	}
	for rowIndex, dataItem := range dataInput {
		fields := slicekit.Keys(dataItem)
		for colIndex, field := range fields {
			cellName := fmt.Sprintf("%s%d", charfileds[colIndex], rowIndex+1)
			value := dataItem[field]
			s.file.SetCellValue(s.sheetname, cellName, value)
		}
	}
	return s.file.WriteToBuffer()
}
