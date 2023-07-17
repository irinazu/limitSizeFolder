package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var root, sor string
	var limit int

	flag.StringVar(&root, "root", "", "root directory")
	flag.IntVar(&limit, "limit", 1024, "limit size")
	flag.StringVar(&sor, "sor", "asc", "sor order")

	flag.Parse()

	moreThenLimitSize, lowThenLimitSize := checkFile(root, limit)

	lowThenLimitSize = checkSlice(lowThenLimitSize, sor)
	moreThenLimitSize = checkSlice(moreThenLimitSize, sor)

	writeDirLowThenLimitSize(lowThenLimitSize)
	writeDirMoreThenLimitSize(moreThenLimitSize)

}

func checkFile(root string, limit int) ([]string, []string) {
	moreThenLimitSize := []string{}
	lowThenLimitSize := []string{}
	var totalSize int64

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			totalSize += info.Size()

			if info.IsDir() {
				path := fmt.Sprintf("%s\t size: %v\n", path, totalSize)
				if totalSize > int64(limit) {
					moreThenLimitSize = append(moreThenLimitSize, path)
				} else {
					lowThenLimitSize = append(lowThenLimitSize, path)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return moreThenLimitSize, lowThenLimitSize
}

// checkSlice() проверяет на пустоту срез из путей и
// переворачивает срез в случае наличия параметра DESC
func checkSlice(slicePathsDir []string, sor string) []string {
	if len(slicePathsDir) == 0 {
		return nil
	}
	if sor == "desc" {
		for i, j := 0, len(slicePathsDir)-1; i < j; i, j = i+1, j-1 {
			slicePathsDir[i], slicePathsDir[j] = slicePathsDir[j], slicePathsDir[i]
		}
	}
	return slicePathsDir
}

// writeDirLowThenLimitSize() выводим в консоль диретории размером меньше чем limit
func writeDirLowThenLimitSize(lowThenLimitSize []string) {
	for _, value := range lowThenLimitSize {
		fmt.Print(value)
	}
}

// writeDirLowThenLimitSize() выводим в файл диретории размером больше чем limit
func writeDirMoreThenLimitSize(moreThenLimitSize []string) {
	fileForOutputDir, errCreateFile := os.Create("output.txt")
	defer fileForOutputDir.Close()
	if errCreateFile != nil {
		log.Println(errCreateFile)
	}

	for _, value := range moreThenLimitSize {
		fileForOutputDir.WriteString(value)
	}
}
