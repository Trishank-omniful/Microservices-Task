package models

import "database/sql"

func ToNullFloat64(val float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: val,
		Valid:   true,
	}
}
