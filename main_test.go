package main

import (
	"fmt"
	"io/fs"
	"os"
	"slices"
	"testing"
)

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

func Test_verifyNumber(t *testing.T) {
	testCases := []struct {
		s string
		n int
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
		// Case 002
		{
			s: "5",
			n: 5,
			m: nil,
		},
		// Case 003
		{
			s: "10",
			n: 10,
			m: nil,
		},
		// Case 004
		{
			s: "50",
			n: 50,
			m: nil,
		},
		// Case 005
		{
			s: "100",
			n: 100,
			m: nil,
		},

		//
		// with %
		//

		// Case 006
		{
			s: "0%",
			n: 0,
			m: nil,
		},
		// Case 007
		{
			s: "5%",
			n: 5,
			m: nil,
		},
		// Case 008
		{
			s: "10%",
			n: 10,
			m: nil,
		},
		// Case 009
		{
			s: "50%",
			n: 50,
			m: nil,
		},
		// Case 010
		{
			s: "100%",
			n: 100,
			m: nil,
		},

		//
		// with spaces
		//

		// Case 006
		{
			s: " 0 ",
			n: 0,
			m: nil,
		},
		// Case 007
		{
			s: "5 %",
			n: 5,
			m: nil,
		},
		// Case 008
		{
			s: "10 %   ",
			n: 10,
			m: nil,
		},
		// Case 009
		{
			s: "   50%",
			n: 50,
			m: nil,
		},
		// Case 010
		{
			s: "   100    %     ",
			n: 100,
			m: nil,
		},

		//
		// errors
		//

		// Case 006
		{
			s: "",
			n: 0,
			m: isStringToNumber,
		},
		// Case 007
		{
			s: "%",
			n: 0,
			m: isStringToNumber,
		},
		// Case 008
		{
			s: " zero ",
			n: 0,
			m: isStringToNumber,
		},
		// Case 009
		{
			s: "-50%",
			n: 0,
			m: isOutOfRange,
		},
		// Case 010
		{
			s: "2354%",
			n: 0,
			m: isOutOfRange,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			n, err := verifyNumber(tc.s)

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

type fakeDirEntry struct {
	name string
}

func (f fakeDirEntry) Name() string               { return f.name }
func (f fakeDirEntry) IsDir() bool                { return false }
func (f fakeDirEntry) Type() fs.FileMode          { return 0 }
func (f fakeDirEntry) Info() (fs.FileInfo, error) { return nil, nil }
