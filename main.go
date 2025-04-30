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
	envPainpointsDataBucket, defPainpointsDataBucket = "PAINPOINTS_DATA_BUCKET", "0,10,20,30,40,50"
	envPainpointsDataDir, defPainpointsDataDir       = "PAINPOINTS_DATA_DIR", "./painpoints/"
	envPainpointsDataSuffix, defPainpointsDataSuffix = "PAINPOINTS_DATA_SUFFIX", "%"
)

const (
	green3 = "\033[38;5;10m\033[48;5;10m"   // brightest green
	green2 = "\033[38;5;112m\033[48;5;112m" //
	green1 = "\033[38;5;220m\033[48;5;220m" // yellow
	red1   = "\033[38;5;214m\033[48;5;214m" // orange
	red2   = "\033[38;5;204m\033[48;5;204m" //
	red3   = "\033[38;5;9m\033[48;5;9m"     // darkest red
)

var colours = []string{
	green3, // >=  0
	green2, // >= 10
	green1, // >= 20
	red1,   // >= 30
	red2,   // >= 40
	red3,   // >= 50
}

var buckets []float64

func init() {
	lis := strings.SplitSeq(dataBucket(), ",")

	for s := range lis {
		n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			tracer.Panic(tracer.Maskf(stringToNumberError, "invalid bucket %q", s))
		}

		{
			buckets = append(buckets, n)
		}
	}

	if len(buckets) != len(colours) {
		tracer.Panic(tracer.Maskf(bucketColourMatchError, "%d != %d", len(buckets), len(colours)))
	}
}

func main() {
	// generate trending percentages, from high to low

	var all []float64
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
				fmt.Printf("%s#%s", colourForNumber(all[i]), "\033[0m")
			} else {
				break
			}
		}

		{
			fmt.Println()
		}
	}
}

func colourForNumber(n float64) string {
	for i := len(buckets) - 1; i >= 0; i-- {
		if n >= buckets[i] {
			return colours[i]
		}
	}

	return colours[0]
}

func dataBucket() string {
	var e = os.Getenv(envPainpointsDataBucket)
	if e != "" {
		return e
	}

	return defPainpointsDataBucket
}

func dataDir() string {
	var e = os.Getenv(envPainpointsDataDir)
	if e != "" {
		return e
	}

	return defPainpointsDataDir
}

func dataSuffix() string {
	var e = os.Getenv(envPainpointsDataSuffix)
	if e != "" {
		return e
	}

	return defPainpointsDataSuffix
}

func fileData() []float64 {
	var all []float64

	for _, p := range readDir() {
		// open the file

		f, err := os.Open(p)
		if err != nil {
			tracer.Panic(err)
		}
		defer f.Close()

		// parse and verify

		n, err := readNumber(readLine(f))
		if err != nil {
			tracer.Panic(err)
		}

		{
			all = append(all, n)
		}
	}

	return all
}

func readDir() []string {
	var pat = dataDir()

	dir, err := os.ReadDir(pat)
	if err != nil {
		tracer.Panic(err)
	}

	// sort dir entries numerically

	return sortPaths(pat, dir)
}

func sortPaths(pat string, dir []os.DirEntry) []string {
	var num []int

	for _, e := range dir {
		if e.Name() == "README.md" {
			continue // ignore data dir readme
		}

		n, err := strconv.Atoi(e.Name())
		if err != nil {
			tracer.Panic(err)
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

func fakeData() []float64 {
	var all []float64

	var tot = 584
	var min = 5
	var max = 60

	for i := range tot {
		var p = float64(i) / float64(tot)
		var b = float64(max)*(1-p) + float64(min)*p
		var n = rand.Float64()*20 - 10 // +-10% noise
		var v = b + n

		{
			all = append(all, v)
		}
	}

	return all
}

func trimString(s string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), dataSuffix()))
}

func readNumber(s string) (float64, error) {
	n, err := strconv.ParseFloat(trimString(s), 64)
	if err != nil {
		return 0, tracer.Maskf(stringToNumberError, "invalid data %q", s)
	}

	return n, nil
}
