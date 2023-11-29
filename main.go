package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

const (
	CIRCLE = "o"
	CROSS  = "x"
)

func main() {
	fmt.Println(Bingo())
}

func Bingo() string {
	bord, lines := makeBord()
	if lines == 0 {
		// 0の場合は無条件Yesという仕様
		return "Yes"
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var result bool
	wg.Add(4)

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
	go commonFn(isBottomRight)
	go commonFn(isBottomLeft)

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
		if cell != CIRCLE && cell != CROSS {
			panic(fmt.Sprintf("Invalid input: Please enter either %s or %s.", CIRCLE, CROSS))
		}
		row = append(row, cell)
	}
	return row
}

func isVerticalBingo(bord []string, n int) bool {
	for line := 0; line < n; line++ {
		if strEvery(bordFilter(bord, func(i int) bool {
			return (i+n)%n == line
		}), func(s string) bool {
			return s == CIRCLE
		}) {
			return true
		}
	}
	return false
}

func isHorizontalBingo(bord []string, n int) bool {
	for line := 0; line < n; line++ {
		if strEvery(bordFilter(bord, func(i int) bool {
			return (n*line)-1 < i && i < (n*(line+1))
		}), func(s string) bool {
			return s == CIRCLE
		}) {
			return true
		}
	}
	return false
}

func isBottomRight(bord []string, n int) bool {
	return strEvery(bordFilter(bord, func(i int) bool {
		if i == 0 {
			return true
		}
		return i%(n+1) == 0
	}), func(s string) bool {
		return s == CIRCLE
	})
}

func isBottomLeft(bord []string, n int) bool {
	return strEvery(bordFilter(bord, func(i int) bool {
		if i == 0 || i == (n*n)-1 {
			return false
		}
		return i%(n-1) == 0
	}), func(s string) bool {
		return s == CIRCLE
	})
}

func bordFilter(list []string, fn func(int) bool) []string {
	result := []string{}
	for i, s := range list {
		if fn(i) {
			result = append(result, s)
		}
	}
	return result
}

func strEvery(list []string, fn func(string) bool) bool {
	for _, ele := range list {
		if !fn(ele) {
			return false
		}
	}
	return true
}
