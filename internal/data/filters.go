package data

import (
	"math"

	"github.com/TH-takahirohara/reading_record_api/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "0より大きい値を指定してください")
	v.Check(f.Page <= 10000, "page", "10000以下の値を指定してください")
	v.Check(f.PageSize > 0, "page_size", "0より大きい値を指定してください")
	v.Check(f.PageSize <= 100, "page_size", "100以下の値を指定してください")
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return f.PageSize * (f.Page - 1)
}

type Metadata struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
