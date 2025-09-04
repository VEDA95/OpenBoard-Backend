package repository

import "gorm.io/gorm"

type QueryOptions struct {
	Select  []string
	Omit    []string
	Preload []string
}

func (options *QueryOptions) AppendToQuery(db *gorm.DB) *gorm.DB {
	query := db

	if len(options.Select) > 0 {
		query = query.Select(options.Select)
	}

	if len(options.Omit) > 0 {
		query = query.Omit(options.Omit...)
	}

	if len(options.Preload) > 0 {
		for _, preloadOption := range options.Preload {
			query = query.Preload(preloadOption)
		}
	}

	return query
}
