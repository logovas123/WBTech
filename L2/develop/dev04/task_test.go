package main

import (
	"reflect"
	"testing"
)

type TestCase struct {
	Dict        *[]string
	ExpectedMap *map[string]*[]string
}

/*
Табличные тесты.
*/

func TestGetCurrentTime(t *testing.T) {
	cases := []TestCase{
		{
			Dict: &[]string{"слиток", "тяпка", "пятак", "листок", "пятка", "столик"},
			ExpectedMap: &map[string]*[]string{
				"слиток": {"листок", "слиток", "столик"},
				"тяпка":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			Dict:        &[]string{"а", "б", "в"},
			ExpectedMap: &map[string]*[]string{},
		},
		{
			Dict:        &[]string{"а", "бб", "бб", "в"},
			ExpectedMap: &map[string]*[]string{},
		},
		{
			Dict:        &[]string{},
			ExpectedMap: &map[string]*[]string{},
		},
	}

	for caseNum, item := range cases {
		resultMap := FindAnagramm(item.Dict)

		if !reflect.DeepEqual(resultMap, item.ExpectedMap) {
			t.Errorf("[%d] wrong results: got %+v, expected %+v", caseNum, resultMap, item.ExpectedMap)
		}

	}
}
