package bookingutils

import (
	"fmt"
	"math"
	"strconv"
)

// StringSliceToInt converts slice of string to slice of uints
func StringSliceToInt(strings []string) (convertedUints []uint, err error) {
	convertedUints = make([]uint, 0, len(strings))
	for _, str := range strings {
		convertedUint, err := strconv.Atoi(str)
		if err != nil {
			return []uint{}, err
		}
		convertedUints = append(convertedUints, uint(convertedUint))
	}
	return
}

func StringToUint(input string) (uint, error) {
	converted, err := strconv.Atoi(input)
	return uint(converted), err
}

func StrToFloat(strValue string) (float64, error) {
	return strconv.ParseFloat(strValue, 64)
}

func UintToString(uintValue uint) string {
	return strconv.FormatUint(uint64(uintValue), 10)
}

func FormatFloatToAmmount(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func AlmostZero(val float64) bool {
	return math.Abs(val) <= 0.05
}
