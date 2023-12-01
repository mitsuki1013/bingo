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

/*
isVerticalBingo
縦列が一つでも一致すればtrueを返す
*/
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

/*
isHorizontalBingo
横列が一つでも一致すればtrueを返す
*/
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

/*
isBottomRight
右下がり斜めが一致すればtrueを返す
*/
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

/*
isBottomLeft
左下り斜めが一致すればtrueを返す
*/
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

/*
bordFilter
縦、横、斜めそれぞれの列を抽出するための抽象関数
第一引数の 'list' の各要素に対して第二引数の関数（fn）を適用する。
また、第二引数は、第一引数のlistの各要素のindexを引数にとる。

ex.
1st ["1", "2", "3"]
2nd func(i int) bool { return i != 2 }
return ["1", "3"]
*/
func bordFilter(list []string, fn func(int) bool) []string {
	result := []string{}
	for i, s := range list {
		if fn(i) {
			result = append(result, s)
		}
	}
	return result
}

/*
strEvery
第一引数（[]string）の各要素に対し、第二引数の関数で評価をした結果、
全てがtrueであればtrueを返し、一つでもfalseがあればfalseを返す。

ex.
1st ["o", "x", "o"]
2nd func(s string) bool { return s === "o" }
return false

ex.
1st ["o", "o", "o"]
2nd func(s string) bool { return s === "o" }
return true
*/
func strEvery(list []string, fn func(string) bool) bool {
	for _, ele := range list {
		if !fn(ele) {
			return false
		}
	}
	return true
}
