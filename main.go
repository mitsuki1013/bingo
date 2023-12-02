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
func isVerticalBingo(bord []string, lines int) bool {
	for i := 0; i < lines; i++ {
		if strEvery(bordFilter(bord, func(index int) bool {
			return (index+lines)%lines == i
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
func isHorizontalBingo(bord []string, lines int) bool {
	for i := 0; i < lines; i++ {
		if strEvery(bordFilter(bord, func(index int) bool {
			return (lines*i)-1 < index && index < (lines*(i+1))
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
func isBottomRight(bord []string, lines int) bool {
	return strEvery(bordFilter(bord, func(index int) bool {
		if index == 0 {
			return true
		}
		return index%(lines+1) == 0
	}), func(s string) bool {
		return s == CIRCLE
	})
}

/*
isBottomLeft
左下り斜めが一致すればtrueを返す
*/
func isBottomLeft(bord []string, lines int) bool {
	return strEvery(bordFilter(bord, func(index int) bool {
		if index == 0 || index == (lines*lines)-1 {
			return false
		}
		return index%(lines-1) == 0
	}), func(s string) bool {
		return s == CIRCLE
	})
}

/*
bordFilter
縦、横、斜めそれぞれの列を抽出するための抽象関数
*/
func bordFilter(bord []string, fn func(int) bool) []string {
	// 1マスのビンゴの場合、縦、横、斜めの全てに置いて条件を満たすため早期リターン
	if len(bord) == 1 {
		return bord
	}
	return filterStrSliceByIndex(bord, fn)
}

/*
filterStrSliceByIndex
第一引数の 'list' の各要素に対して第二引数の関数（fn）を適用する。
また、第二引数の関数は、第一引数の 'list' の各要素のindexを引数にとる。

ex.
1st ["1", "2", "3"]
2nd func(i int) bool { return i != 2 }
return ["1", "3"]
*/
func filterStrSliceByIndex(list []string, fn func(int) bool) []string {
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
	for _, s := range list {
		if !fn(s) {
			return false
		}
	}
	return true
}
