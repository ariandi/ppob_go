package util

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

func IntToStr(val int) string {
	return strconv.Itoa(val)
}

func Int64ToStr(val int64) string {
	return strconv.FormatInt(val, 10)
}

func StrToInt(val string) (int, error) {
	s := val
	if n, err := strconv.Atoi(s); err == nil {
		return n, nil
	} else {
		logrus.Info("[Util Convert] err when convert string to int : ", err)
		return 0, err
	}
}

func StrToInt64(val string) (int64, error) {
	s := val
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		logrus.Info("[Util Convert] err when convert string to int64 : ", err)
		return 0, err
	}

	return n, nil
}
