package main

import (
	"reflect"
	"testing"
)

type TestCase struct {
	ExtractString  string
	ExpectedString string
}

/*
Табличные тесты.
*/

func TestGetCurrentTime(t *testing.T) {
	cases := []TestCase{
		{ExtractString: "a4bc2d5e", ExpectedString: "aaaabccddddde"},
		{ExtractString: "a10bc2d5e", ExpectedString: "aaaaaaaaaabccddddde"},
		{ExtractString: "abcd", ExpectedString: "abcd"},
		{ExtractString: "abcd5", ExpectedString: "abcddddd"},
		{ExtractString: "a4bc2d5e10", ExpectedString: "aaaabccdddddeeeeeeeeee"},
		{ExtractString: "a", ExpectedString: "a"},
		{ExtractString: "a2", ExpectedString: "aa"},
		{ExtractString: "a10", ExpectedString: "aaaaaaaaaa"},
		{ExtractString: "45", ExpectedString: ""},
		{ExtractString: "", ExpectedString: ""},
	}

	for caseNum, item := range cases {
		unpackString := getUnpackString(item.ExtractString)

		if !reflect.DeepEqual(unpackString, item.ExpectedString) {
			t.Errorf("[%d] wrong results: got %+v, expected %+v", caseNum, unpackString, item.ExpectedString)
		}

	}
}
