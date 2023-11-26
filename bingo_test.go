package main

import (
	"os"
	"testing"
)

func TestBingo(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe(): %v", err)
	}

	// テスト用にos.Stdin、os.StdoutをReader、Writerと差し替え
	osStdin, osStdout := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = r, w
	defer func() { os.Stdin, os.Stdout = osStdin, osStdout }()

	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "0/一致",
			in:   "0\n",
			want: "Yes",
		},
		{
			name: "1/一致",
			in: "1\n" +
				"o\n",
			want: "Yes",
		},
		{
			name: "2/横一致",
			in: "2\n" +
				"oo\n" +
				"xx\n",
			want: "Yes",
		},
		{
			name: "3/横一致",
			in: "3\n" +
				"ooo\n" +
				"oxx\n" +
				"xox\n",
			want: "Yes",
		},
		{
			name: "3/縦一致",
			in: "3\n" +
				"xoo\n" +
				"oox\n" +
				"xox\n",
			want: "Yes",
		},
		{
			name: "3/右下がり斜め一致",
			in: "3\n" +
				"oox\n" +
				"oox\n" +
				"xoo\n",
			want: "Yes",
		},
		{
			name: "3/左下がり斜め一致",
			in: "3\n" +
				"oxo\n" +
				"xox\n" +
				"oox\n",
			want: "Yes",
		},
		{
			name: "3/不一致",
			in: "3\n" +
				"oxo\n" +
				"xox\n" +
				"xox\n",
			want: "No",
		},
		{
			name: "5/不一致",
			in: "5\n" +
				"xoooo\n" +
				"oxooo\n" +
				"ooxoo\n" +
				"oooxo\n" +
				"oooox\n",
			want: "No",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input := []byte(c.in)

			if n, err := os.Stdout.Write(input); err != nil {
				t.Errorf("input is %v bytes, but only %v byte written", len(input), n)
				return
			}

			got := Bingo()
			if got != c.want {
				t.Errorf("failed: expected %v, got %v", c.want, got)
				return
			}
		})
	}
}
