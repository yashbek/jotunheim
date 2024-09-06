package utils

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}

func Cmp(a, b int) bool {
	return a - b >= 0
}

func WithinBounds(num, inc, exc int) bool {
	return num >= inc && num < exc
}