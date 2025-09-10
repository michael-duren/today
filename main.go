package main

import (
	"fmt"

	todayfile "today/today-file"
)

func main() {
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
