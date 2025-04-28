package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"slices"

	"github.com/xh3b4sd/tracer"
)

const (
	envHeatmapDataDir, defHeatmapDataDir = "HEATMAP_DATA_DIR", "./manual-testing/"
)

// Buckets for color mapping
var colors = []string{
	"\033[48;5;9m",   // Red
	"\033[48;5;204m", //
	"\033[48;5;214m", //
	"\033[48;5;220m", //
	"\033[48;5;112m", //
	"\033[48;5;10m",  // Green
}

func main() {
	// generate trending percentages, from high to low

	var all []int
	if len(os.Args) == 2 && (os.Args[1] == "-t" || os.Args[1] == "--test") {
		all = fakeData()
	} else {
		all = fileData()
	}

	// print heatmap, from left to right

	var row = 7
	var col = (len(all) / row) + 1

	for r := range row {
		for c := range col {
			var i = (c * row) + r

			if i < len(all) {
				fmt.Printf("%s %s", colorForPercentage(all[i]), "\033[0m")
			} else {
				break
			}
		}

		{
			fmt.Println()
		}
	}
}

func colorForPercentage(p int) string {
	switch {
	case p > 50:
		return colors[0]
	case p > 40:
		return colors[1]
	case p > 30:
		return colors[2]
	case p > 20:
		return colors[3]
	case p > 10:
		return colors[4]
	default:
		return colors[5]
	}
}

func dataPath() string {
	var e = os.Getenv(envHeatmapDataDir)
	if e != "" {
		return e
	}

	return defHeatmapDataDir
}

func fileData() []int {
	var all []int

	for _, p := range readDir() {
		// open the file

		f, err := os.Open(p)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// parse and verify

		n, err := verifyNumber(readLine(f))
		if err != nil {
			panic(err)
		}

		{
			all = append(all, int(n))
		}
	}

	return all
}

func minMax(f float64) int {
	if f > 100 {
		return 100
	} else if f < 0 {
		return 0
	}

	return int(f)
}

func readDir() []string {
	var pat = dataPath()

	dir, err := os.ReadDir(pat)
	if err != nil {
		panic(err)
	}

	// sort dir entries numerically

	return sortPaths(pat, dir)
}

func sortPaths(pat string, dir []os.DirEntry) []string {
	var num []int

	for _, e := range dir {
		n, err := strconv.Atoi(e.Name())
		if err != nil {
			panic(err)
		}

		{
			num = append(num, n)
		}
	}

	{
		slices.Sort(num)
	}

	var str []string
	for _, n := range num {
		str = append(str, filepath.Join(pat, strconv.Itoa(n)))
	}

	return str
}

func readLine(f *os.File) string {
	s := bufio.NewScanner(f)

	if s.Scan() {
		return s.Text()
	}

	return ""
}

func fakeData() []int {
	var all []int

	var tot = 584
	var min = 5
	var max = 60

	for i := range tot {
		var p = float64(i) / float64(tot)
		var b = float64(max)*(1-p) + float64(min)*p
		var n = rand.Float64()*20 - 10 // +-10% noise
		var v = b + n

		{
			all = append(all, minMax(v))
		}
	}

	return all
}

func trimString(s string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "%"))
}

func verifyNumber(s string) (int, error) {
	n, err := strconv.ParseFloat(trimString(s), 64)
	if err != nil {
		return 0, tracer.Maskf(stringToNumberError, "%s", s)
	}

	if n < 0 || n > 100 {
		return 0, tracer.Maskf(outOfRangeError, "%s", s)
	}

	return int(n), nil
}
