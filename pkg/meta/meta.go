package meta

import (
	"os"
	"strconv"
)

type Meta struct {
	Page       int `json:"page"`
	SizePage   int `json:"size_page"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}

func New(page, sizePage, total int) (*Meta, error) {
	if sizePage <= 0 {
		var err error
		sizePage, err = strconv.Atoi(os.Getenv("PAGINATOR_LIMIT_DEFAULT"))
		if err != nil {
			return nil, err
		}
	}
	pageCount := 0

	if total >= 0 {
		pageCount = (total + sizePage - 1) / sizePage
		if page >= pageCount {
			page = pageCount
		}
		if page < 1 {
			page = 1
		}
	}
	return &Meta{
		Page:       page,
		SizePage:   sizePage,
		PageCount:  pageCount,
		TotalCount: total,
	}, nil
}

func (p *Meta) Offser() int {
	return (p.Page - 1) * p.SizePage
}

func (p *Meta) Limit() int {
	return p.SizePage
}
