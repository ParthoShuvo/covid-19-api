package resource

import (
	"net/url"
	"strconv"

	"github.io/covid-19-api/errors"
)

type pageInfo struct {
	page     int
	pageSize int
}

func parsePageInfo(query url.Values) (*pageInfo, error) {
	var (
		page            int
		pageSize        int
		defaultPage     int = 1
		defaultPageSize int = 10
	)
	anyOrElse := func(val string, defaultVal int) int {
		if parsedVal, err := strconv.Atoi(val); err == nil {
			return parsedVal
		}
		return defaultVal
	}
	isNeg := func(val int) bool { return val < 0 }
	page, pageSize = anyOrElse(query.Get("page"), defaultPage), anyOrElse(query.Get("pagesize"), defaultPageSize)
	if isNeg(page) || isNeg(pageSize) {
		return &pageInfo{}, errors.BadRequest.New("malformed query param values in page & pagesize")
	}
	return &pageInfo{page, pageSize}, nil
}
