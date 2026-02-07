package util

import "errors"

var widthOfA4 = 210

func MmToPixel(num, pageWidth int) int {
	return int(float64(num) / float64(widthOfA4) * float64(pageWidth))
}

func CheckNonNegative(nums ...int) error {
	for _, num := range nums {
		if num < 0 {
			return errors.New("negative flag values are invalid")
		}
	}
	return nil
}
