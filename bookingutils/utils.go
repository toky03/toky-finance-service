package bookingutils

import "strconv"

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
