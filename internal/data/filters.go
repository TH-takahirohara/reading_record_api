package data

import "github.com/TH-takahirohara/reading_record_api/internal/validator"

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

func (f Filters) offset() int {
	return f.PageSize * (f.Page - 1)
}
