package actions

import (
	"errors"
	"fmt"
	"slices"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Trace struct {
	key        string
	value      interface{}
	oldValue   interface{}
	changeType string
	file       string
}

func TracesToString(traces *[]Trace) string {
	var result string
	for _, trace := range *traces {
		result += trace.changeType
		result += " "
		result += trace.key
		result += " "
		result += fmt.Sprint("Value: ", trace.value)
		if trace.oldValue != nil {
			result += " "
			result += fmt.Sprint("Old Value: ", trace.oldValue)
		}
		result += " "
		result += trace.file
		result += "\n"
	}
	return result
}

func traceToTable(traces *[]Trace, configs []Config, fileType string, memoryMap map[string]interface{}) (string, error) {

	columns := []string{}
	t := table.NewWriter()
	// t.SetOutputMirror(os.Stdout)
	header := table.Row{""}
	for order, config := range configs {
		header = append(header, fmt.Sprintf("%s-%d", config.Path, order))
		columns = append(columns, fmt.Sprintf("%s-%d", config.Path, order))
	}

	header = append(header, "PreTemplate")
	header = append(header, "PostTemplate")
	columns = append(columns, "PreTemplate")
	columns = append(columns, "PostTemplate")

	t.AppendHeader(header)
	keys := findKeys(traces)
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	slices.Sort(sortedKeys)
	sortedKeys = slices.Compact(sortedKeys)
	info := buildTableInfo(traces, columns, sortedKeys, memoryMap)
	buildTableBody(info, columns, t, sortedKeys)
	t.SetStyle(table.Style{
		Name: "trace",
		Color: table.ColorOptions{
			Row:          text.Colors{text.BgHiCyan, text.BgBlack},
			RowAlternate: text.Colors{text.BgCyan, text.FgBlack}},
	})

	switch fileType {
	case "md":
		return t.RenderMarkdown(), nil
	case "html":
		return t.RenderHTML(), nil
	case "csv":
		return t.RenderCSV(), nil
	case "txt":
		return t.Render(), nil
	default:
		return t.Render(), errors.New("file type not supported for trace output (supported types: .md, .html, .csv, .txt)")
	}

}

func buildTableInfo(traces *[]Trace, columns []string, keys []string, memoryMap map[string]interface{}) Table {

	rows := make(map[string]Row)
	for _, key := range keys {
		rows[key] = buildRowInfo(traces, columns, key, memoryMap[key])
	}

	return Table{Rows: rows}

}

func buildRowInfo(traces *[]Trace, columns []string, key string, finalValue interface{}) Row {
	row := Row{Columns: make(map[string]Metadata)}

	for _, trace := range *traces {
		if trace.key == key {
			for _, column := range columns {
				if trace.file == column {
					row.Columns[column] = Metadata{value: trace.value, crud: trace.changeType}
					row.Columns["PreTemplate"] = Metadata{value: trace.value, crud: "none"}
				}
			}
		}
	}
	row.Columns["PostTemplate"] = Metadata{value: finalValue, crud: "none"}
	return row
}

func findKeys(traces *[]Trace) []string {
	keys := make([]string, len(*traces))
	for i, trace := range *traces {
		keys[i] = trace.key
	}
	return keys
}

func buildTableBody(info Table, columns []string, t table.Writer, keys []string) {
	for _, key := range keys {
		row := info.Rows[key]
		rowData := []interface{}{key}
		for _, column := range columns {
			if row.Columns[column].crud == "delete" {
				rowData = append(rowData, "(deleted)")
				continue
			}
			if row.Columns[column].value == nil {
				rowData = append(rowData, "")
				continue
			}
			rowData = append(rowData, row.Columns[column].value)
		}
		t.AppendRow(rowData)
	}

}

type Table struct {
	Rows map[string]Row
}

type Row struct {
	Columns map[string]Metadata
}

type Metadata struct {
	value interface{}
	crud  string
}
