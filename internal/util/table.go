package util

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintTableWithAny(arr any) (err error) {

	if bts, err := json.Marshal(&arr); err != nil {
		return err
	} else {
		var result []map[string]any = make([]map[string]any, 0)
		if err = json.Unmarshal(bts, &result); err != nil {
			return err
		} else {
			return PrintTable(result)
		}
	}
}

func PrintTableWithArray(arr []any) (err error) {
	var result []map[string]any = make([]map[string]any, 0)
	if bts, err := json.Marshal(&arr); err != nil {
		return err
	} else {
		if err = json.Unmarshal(bts, &result); err != nil {
			return err
		} else {
			return PrintTable(result)
		}
	}
}

func PrintTableWithByteArray(bts []byte) (err error) {
	var result []map[string]interface{} = make([]map[string]interface{}, 0)
	if err = json.Unmarshal(bts, &result); err != nil {
		return err
	} else {
		return PrintTable(result)

	}
}

func guessTitle(rows []map[string]interface{}) []interface{} {
	if len(rows) == 0 {
		return make([]interface{}, 0)
	}
	titleArr := make([]interface{}, 0)
	for k, _ := range rows[0] {
		titleArr = append(titleArr, k)
	}
	return titleArr
}
func PrintTable(rows []map[string]interface{}) (err error) {
	titleArr := guessTitle(rows)
	writer := table.NewWriter()
	header := table.Row{}
	header = append(header, titleArr...)
	//header = append(header, titleArr...)
	writer.AppendHeader(header)
	writer.SetAutoIndex(true)
	writer.Style().Options.SeparateRows = true
	for j, _ := range rows {
		row := rows[j]
		obj := table.Row{}
		for i, _ := range titleArr {
			title := titleArr[i]
			strtitle := title.(string)
			obj = append(obj, row[strtitle])
		}
		writer.AppendRow(obj)
	}
	fmt.Println(writer.Render())
	return nil
}
