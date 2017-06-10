package redis

import "errors"

//GetPageRangeForStrings returns offset start and end for given page
func GetPageRangeForStrings(values []string, pageSize int, currentPage int) (int, int, error) {

	pageOffsetStart := (currentPage - 1) * pageSize
	if pageOffsetStart > len(values) {
		return 0, 0, errors.New("page not found")
	}

	pageOffsetEnd := currentPage * pageSize
	if pageOffsetEnd > len(values) {
		pageOffsetEnd = len(values)
	}

	return pageOffsetStart, pageOffsetEnd, nil
}
