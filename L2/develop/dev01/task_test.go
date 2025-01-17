package main

import (
	"reflect"
	"testing"
	"time"
)

type TestCase struct {
	NtpServer  string
	ExpectTime time.Time
	IsError    bool
}

/*
Табличные тесты:
Два теста: тест с правильным рез-м, и тест с ошибкой (ошибка имени сервера).
В первом тесте проверяю что время сервера и текущее время отличается менее чем на одну секунду.
Вовтором проверяю, что вернулась пустая структура.
*/

func TestGetCurrentTime(t *testing.T) {
	cases := []TestCase{
		{NtpServer: "0.ru.pool.ntp.org", IsError: false},
		{NtpServer: "bad server", ExpectTime: time.Time{}, IsError: true},
	}

	for caseNum, item := range cases {
		ntpTime, err := getCurrentTime(item.NtpServer)

		if item.IsError && err == nil {
			t.Errorf("[%v] expected error, got nil", caseNum)
		}

		if !item.IsError && err != nil {
			t.Errorf("[%d] unexpected error: %v", caseNum, err)
		}

		switch caseNum {
		case 0:
			if time.Since(ntpTime) > time.Second {
				t.Errorf("[%d] wrong results: got bad time from ntp server", caseNum)
			}
		default:
			if !reflect.DeepEqual(ntpTime, item.ExpectTime) {
				t.Errorf("[%d] wrong results: got %+v, expected %+v", caseNum, ntpTime, item.ExpectTime)
			}
		}
	}
}
