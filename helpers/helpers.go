package helpers

import "strconv"

//SizeInBytesToHumanReadable converts bytes count to human-readable string
//e.g. SizeInBytesToHumanReadable(0) = 0B
//SizeInBytesToHumanReadable(10240) = 10K
func SizeInBytesToHumanReadable(size int64) string {
	values := []string{"K", "MB", "GB", "TB"}

	humanSize := strconv.Itoa(int(size)) + "B"
	i := 0
	for size/1024 > 0 && i < len(values) {
		size = size / 1024
		humanSize = strconv.Itoa(int(size)) + values[i]
		i = i + 1
	}

	return humanSize
}
