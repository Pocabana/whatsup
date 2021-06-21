package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type GameStats struct {
	Total   	string
	Under5GP 	string
	Under5REM	string
	Over5GP 	string
	Over5REM 	string
	Name 		string
	Rank		int
}

func getFilesList(filePath string, csvFiles *[]string){
	var files []string

	//Get the files in the folder
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	//If it is a csv_import file, keep it
	for _, file := range files {
		fileExtension := filepath.Ext(file)
		if fileExtension == ".csv" {
			*csvFiles = append(*csvFiles, file)
		}
	}
}

func main() {
	//Variables
	var filePath string
	var csvFiles []string

	//Read the argument -path to get the csv_import folder
	flag.StringVar(&filePath, "path", "", "Usage")
	flag.Parse()

	//Get the list of CSV files in the folder
	getFilesList(filePath, &csvFiles)

	//Open the csv_import file
	for _, file := range csvFiles {
		csvOpened, err := os.Open(file)
		if err != nil {
			fmt.Println("Cannot open file", file, err)
		}
		defer csvOpened.Close()

		csvLines, err := csv.NewReader(csvOpened).ReadAll()
		if err != nil {
			fmt.Println(err)
		}

		var iteration = 0
		var totalWeekly = 0
		var data []GameStats
		for _, line := range csvLines {

			row := GameStats{
				Total:    	line[0],
				Under5GP: 	line[1],
				Under5REM:	line[2],
				Over5GP:  	line[3],
				Over5REM: 	line[4],
				Name:   	line[5],
				Rank:		iteration,
			}
			data = append(data, row)
			if iteration != 0 {
				totalInt, err := strconv.Atoi(row.Total)
				if err != nil {
					fmt.Println(err)
				}
				if iteration >= 10 {
					totalWeekly += totalInt
				}
			}
			iteration++
		}

		weeklyTopGames := []string{"Top1", "Top2", "Top3", "Top4", "Top5", "Total"}


		csvWeeklyTop, err := os.OpenFile("csv_internal/weeklytop5.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("Cannot open file", file, err)
		}
		for i := 1;  i<=5; i++ {
			weeklyTopGames[i-1] = "" + data[i].Total
		}
		weeklyTopGames[5] = strconv.Itoa(totalWeekly)

		csvWriter := csv.NewWriter(csvWeeklyTop)
		csvWriter.Write(weeklyTopGames)
		csvWriter.Flush()
		csvWeeklyTop.Close()

		//for _, row := range data {
		//	fmt.Println(  row.Rank, "\t" + row.Name + "\t\t\t\t " + row.Total + " \t\t ")
		//}


		fmt.Println("Total weekly: ", totalWeekly)
		fmt.Println(weeklyTopGames)

	}
}
