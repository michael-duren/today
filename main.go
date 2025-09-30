package main

import (
	"fmt"
	"os"
	"strings"

	todayfile "today/today-file"
)

func main() {
	currPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("get cwd failed %s", err.Error()))
	}
	if strings.Contains(currPath, "daily") {
		os.Chdir("daily")
	}

	today := todayfile.NewTodayFile()
	f, err := today.CreateFile()
	if err != nil {
		fmt.Println("unable to create file")
	}
	err = today.UpdateContents(f)
	if err != nil {
		fmt.Println("unable to update file")
		return
	}

	fmt.Println("successfully created file")
}
