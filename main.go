package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

const (
	o = "o"
	x = "x"
)

func main() {
	fmt.Println(Bingo())
}

func Bingo() string {
	bord, lines := makeBord()
	// if lines == 0 || lines == 1 || lines == 2 {
	// 	return "NO"
	// }

	var mu sync.Mutex
	var wg sync.WaitGroup
	var result bool
	wg.Add(3)

	commonFn := func(fn func([]string, int) bool) {
		defer wg.Done()
		if found := fn(bord, lines); found {
			mu.Lock()
			result = true
			mu.Unlock()
		}
	}

	go commonFn(isVerticalBingo)
	go commonFn(isHorizontalBingo)
	go commonFn(isDiagonalBingo)

	wg.Wait()
	if result {
		return "Yes"
	}
	return "No"
}

func makeBord() ([]string, int) {
	scanner := bufio.NewScanner(os.Stdin)
	var bord []string

	scanner.Scan()
	lines, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	for i := 0; i < lines+1; i++ {
		if i == 0 {
			continue
		}
		scanner.Scan()
		row := makeRow(scanner.Text())
		if len(row) != lines {
			panic(fmt.Sprintf("Invalid input: The number of input cells must be exactly %d.", lines))
		}
		bord = append(bord, row...)
	}

	return bord, lines
}

func makeRow(inputRow string) []string {
	var row []string
	for _, char := range inputRow {
		cell := string(char)
		if cell != o && cell != x {
			panic(`Invalid input: Please enter either "o" or "x".`)
		}
		row = append(row, cell)
	}
	return row
}

func isVerticalBingo(bord []string, lines int) bool {
	for columnNo := 0; columnNo < lines; columnNo++ {
		upperColumn := bord[columnNo]
		result := true
		if upperColumn == x {
			result = false
			continue
		}
		for rowNo := 0; rowNo < lines; rowNo++ {
			if rowNo == 0 {
				continue
			}
			if upperColumn != bord[(lines*rowNo)+columnNo] {
				result = false
			}
		}
		if result {
			return true
		}
	}
	return false
}

func isHorizontalBingo(bord []string, lines int) bool {
	for rowNo := 0; rowNo < lines; rowNo++ {
		firstColumnNo := lines * rowNo
		lastColumnNo := firstColumnNo + lines
		if bord[firstColumnNo] == x {
			continue
		}
		if lines == 1 {
			return true
		}
		ts := bord[firstColumnNo+1 : lastColumnNo]
		//fmt.Println(ts)
		if strEvery(ts, func(cell string) bool {
			return bord[firstColumnNo] == cell
		}) {
			return true
		}
	}
	return false
}

func strEvery(list []string, fn func(string) bool) bool {
	for _, ele := range list {
		if !fn(ele) {
			return false
		}
	}
	return true
}

func isDiagonalBingo(bord []string, lines int) bool {
	downRightResult := true
	downLeftResult := true

	if lines == 1 && bord[0] == x {
		return false
	}
	for rowNo := 0; rowNo < lines; rowNo++ {
		upperLastIndex := lines - 1
		firstUpperColumn := bord[0]
		lastUpperColumn := bord[upperLastIndex]
		if rowNo == 0 {
			continue
		}
		if firstUpperColumn == x || firstUpperColumn != bord[(lines+1)*rowNo] {
			downRightResult = false
		}
		if lastUpperColumn == x || lastUpperColumn != bord[upperLastIndex+(rowNo*upperLastIndex)] {
			downLeftResult = false
		}
	}
	return downRightResult || downLeftResult
}
