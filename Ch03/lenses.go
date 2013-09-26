package main

import (
	"./trees"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("lenses.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	lenses := make([]*trees.Point, 0)
	reader := bufio.NewReader(file)
	for {
		line, lineError := reader.ReadString('\n')
		if lineError == io.EOF {
			break
		}
		fields := strings.Split(line, "\t")
		// have to explict convert to []interface{}, otherweise will raise error:
		// cannot use fields (type []string) as type []interface {} in assignment
		ff := make([]interface{}, len(fields))
		for i, field := range fields {
			ff[i] = field
		}
		lense := trees.NewPoint(ff...)
		lenses = append(lenses, lense)
	}
	lensesLabels := []string{"age", "prescript", "astigmatic", "tearRate"}
	lenseTree := trees.CreateTree(lenses, lensesLabels)
	fmt.Println(lenseTree)
}
