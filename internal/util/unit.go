package util

var widthOfA4 = 210

func MmToPixel(num, pageWidth int) int {
	return int(float64(num) / float64(widthOfA4) * float64(pageWidth))
}
