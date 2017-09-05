package helpers

import "strconv"

func SizeInBytesToHumanReadable(size int64) string {
	values := []string{"MB", "GB", "TB"}

	humanSize := strconv.Itoa(int(size)) + "K"
	i := 0
	for size/1024 > 0 && i < len(values) {
		size = size / 1024
		humanSize = strconv.Itoa(int(size)) + values[i]
		i = i + 1
	}

	return humanSize
}
