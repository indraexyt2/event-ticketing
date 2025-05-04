package utils

import "strconv"

type PaginationParams struct {
	Page  int
	Limit int
}

func NewPaginationParams(page, limit string) PaginationParams {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	if limitInt > 100 {
		limitInt = 100
	}

	return PaginationParams{
		Page:  pageInt,
		Limit: limitInt,
	}
}

func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p PaginationParams) GetLimit() int {
	return p.Limit
}
