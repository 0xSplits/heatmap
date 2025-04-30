package main

import (
	"fmt"
	"io/fs"
	"os"
	"slices"
	"testing"
)

func Test_colourForNumber(t *testing.T) {
	testCases := []struct {
		n float64
		c string
	}{
		// Case 000
		{
			n: 0,
			c: green3,
		},
		// Case 001
		{
			n: 7.5,
			c: green3,
		},
		// Case 002
		{
			n: 10,
			c: green2,
		},
		// Case 003
		{
			n: 13,
			c: green2,
		},
		// Case 004
		{
			n: 21,
			c: green1,
		},
		// Case 005
		{
			n: 27.003,
			c: green1,
		},
		// Case 006
		{
			n: 30,
			c: red1,
		},
		// Case 007
		{
			n: 30.1,
			c: red1,
		},
		// Case 008
		{
			n: 39.9,
			c: red1,
		},
		// Case 009
		{
			n: 43,
			c: red2,
		},
		// Case 010
		{
			n: 47,
			c: red2,
		},
		// Case 011
		{
			n: 50.5,
			c: red3,
		},
		// Case 012
		{
			n: 55,
			c: red3,
		},
		// Case 013
		{
			n: 68,
			c: red3,
		},
		// Case 014
		{
			n: 128,
			c: red3,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			c := colourForNumber(tc.n)

			if c != tc.c {
				t.Fatalf("expected %#v got %#v", tc.c, c)
			}
		})
	}
}

func Test_readNumber(t *testing.T) {
	testCases := []struct {
		s string
		n float64
		m func(error) bool
	}{
		//
		// without %
		//

		// Case 000
		{
			s: "0",
			n: 0,
			m: nil,
		},
		// Case 001
		{
			s: "5",
			n: 5,
			m: nil,
		},
		// Case 002
		{
			s: "10",
			n: 10,
			m: nil,
		},
		// Case 003
		{
			s: "-500",
			n: -500,
			m: nil,
		},
		// Case 004
		{
			s: "100",
			n: 100,
			m: nil,
		},

		//
		// with %
		//

		// Case 005
		{
			s: "0%",
			n: 0,
			m: nil,
		},
		// Case 006
		{
			s: "5%",
			n: 5,
			m: nil,
		},
		// Case 007
		{
			s: "10%",
			n: 10,
			m: nil,
		},
		// Case 008
		{
			s: "50%",
			n: 50,
			m: nil,
		},
		// Case 009
		{
			s: "100%",
			n: 100,
			m: nil,
		},

		//
		// with floats
		//

		// Case 010
		{
			s: "0.5%",
			n: 0.5,
			m: nil,
		},
		// Case 011
		{
			s: "0.76",
			n: 0.76,
			m: nil,
		},
		// Case 012
		{
			s: "10.007%",
			n: 10.007,
			m: nil,
		},
		// Case 013
		{
			s: "50.12004",
			n: 50.12004,
			m: nil,
		},
		// Case 014
		{
			s: "-1.5",
			n: -1.5,
			m: nil,
		},

		//
		// with spaces
		//

		// Case 015
		{
			s: " 0 ",
			n: 0,
			m: nil,
		},
		// Case 016
		{
			s: "-5 %",
			n: -5,
			m: nil,
		},
		// Case 017
		{
			s: "10 %   ",
			n: 10,
			m: nil,
		},
		// Case 018
		{
			s: "   50%",
			n: 50,
			m: nil,
		},
		// Case 019
		{
			s: "   150    %     ",
			n: 150,
			m: nil,
		},

		//
		// errors
		//

		// Case 020
		{
			s: "",
			n: 0,
			m: isStringToNumber,
		},
		// Case 021
		{
			s: "%",
			n: 0,
			m: isStringToNumber,
		},
		// Case 022
		{
			s: " zero ",
			n: 0,
			m: isStringToNumber,
		},
		// Case 023
		{
			s: "     ",
			n: 0,
			m: isStringToNumber,
		},
		// Case 024
		{
			s: "ABC",
			n: 0,
			m: isStringToNumber,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			n, err := readNumber(tc.s)

			if tc.m == nil && err != nil {
				t.Fatalf("error without matcher %#v", err)
			} else if tc.m != nil && !tc.m(err) {
				t.Fatalf("matcher returned false %#v", err)
			}

			if n != tc.n {
				t.Fatalf("expected %#v got %#v", tc.n, n)
			}
		})
	}
}

func Test_sortPaths(t *testing.T) {
	dir := []os.DirEntry{
		fakeDirEntry{name: "22"},
		fakeDirEntry{name: "0"},
		fakeDirEntry{name: "7"},
		fakeDirEntry{name: "106"},
		fakeDirEntry{name: "277"},
		fakeDirEntry{name: "11"},
		fakeDirEntry{name: "68"},
		fakeDirEntry{name: "103"},
		fakeDirEntry{name: "1"},
	}

	act := sortPaths("foo", dir)

	exp := []string{
		"foo/0",
		"foo/1",
		"foo/7",
		"foo/11",
		"foo/22",
		"foo/68",
		"foo/103",
		"foo/106",
		"foo/277",
	}

	if !slices.Equal(act, exp) {
		t.Fatalf("expected %#v got %#v", exp, act)
	}
}

type fakeDirEntry struct {
	name string
}

func (f fakeDirEntry) Name() string               { return f.name }
func (f fakeDirEntry) IsDir() bool                { return false }
func (f fakeDirEntry) Type() fs.FileMode          { return 0 }
func (f fakeDirEntry) Info() (fs.FileInfo, error) { return nil, nil }
