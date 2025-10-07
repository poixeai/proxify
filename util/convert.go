package util

import "strconv"

// int 转换成 string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// string 转换为 int
func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	return i, err
}

// int64 转换为 string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// string 转换为 int64
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// float64 转换为 string
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// string 转换为 float64
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// bool 转换成 string
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// string 转换成 bool
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
