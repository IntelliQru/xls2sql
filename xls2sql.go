package main

import (
	"os"
	"fmt"
	"github.com/tealeg/xlsx"
	"xls2sql/config"
	"strings"
)

var cfg = config.NewConfig()

func main() {

	cfg.LoadConfig("config.cfg")


	inputFile := cfg.Get("input")
	if inputFile == "" {
		panic("Invalid config value: input")
	}

	outputFile := cfg.Get("output")
	if outputFile == "" {
		panic("Invalid config value: output")
	}

	xlsFile, err := xlsx.OpenFile(inputFile)
	if err != nil {
		panic(err)
	}

	outFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0777)
	outFile.Truncate(0)
	defer outFile.Close()
	if err != nil {
		panic(err)
	}

	defType := cfg.Get("defaultType")
	if defType == "" {
		defType = "VARCHAR(200)"
	}

	var sql string

	for _, sheet := range xlsFile.Sheets {

		fmt.Print("Sheet " + sheet.Name + "...")

		if len(sheet.Rows) < 1 {
			fmt.Println("No rows")
			continue
		}

		colNames := make([]string, 0)

		sql = "CREATE TABLE IF NOT EXISTS `" + sheet.Name + "` (\n"

		bFirst := true

		bHeader := cfg.Get("headerRow") == "true"
		colPrefix := cfg.Get("colPrefix")

		if colPrefix == "" {
			colPrefix = "col"
		}

		if bHeader {
			for _, cell := range sheet.Rows[0].Cells {
				if !bFirst {
					sql += ", "
				} else {
					bFirst = false
				}

				sql += "`" + cell.Value + "` " + defType
				colNames = append(colNames, cell.Value)
			}
		} else {
			for i := 0; i < len(sheet.Rows[0].Cells); i++ {
				if !bFirst {
					sql += ", "
				} else {
					bFirst = false
				}

				colName := colPrefix + fmt.Sprintf("%d", i + 1)

				sql += colName + " " + defType
				colNames = append(colNames, colName)
			}
		}

		sql += ");\n\n"

		outFile.WriteString(sql)

		totalRows := 0

		for _, row := range sheet.Rows {

			bFirst = true
			sql = "INSERT INTO `" + sheet.Name + "`("
			sqlVals := " VALUES("
			for iCol, colName := range colNames {
				if !bFirst {
					sql += ", "
					sqlVals += ", "
				} else {
					bFirst = false
				}

				sql += "`" + colName + "`"
				if iCol < len(row.Cells) {
					val := strings.Replace(row.Cells[iCol].Value, "'", "\\'", -1)
					sqlVals += "'" + val + "'"
				} else {
					sqlVals += "''"
				}
			}

			sql += ")\n" + sqlVals + ");\n\n"
			outFile.WriteString(sql)
			totalRows++
		}

		fmt.Println(fmt.Sprintf("total %d rows", totalRows))
	}

}