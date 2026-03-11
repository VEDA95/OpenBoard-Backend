package repository

import "gorm.io/gorm"

type queryConfig struct {
	selects  []string
	omits    []string
	preloads []string
}

type QueryOption func(*queryConfig)

func WithSelect(fields ...string) QueryOption {
	return func(c *queryConfig) {
		c.selects = append(c.selects, fields...)
	}
}

func WithOmit(fields ...string) QueryOption {
	return func(c *queryConfig) {
		c.omits = append(c.omits, fields...)
	}
}

func WithPreload(relations ...string) QueryOption {
	return func(c *queryConfig) {
		c.preloads = append(c.preloads, relations...)
	}
}

func applyOptions(db *gorm.DB, opts []QueryOption) *gorm.DB {
	cfg := &queryConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	if len(cfg.selects) > 0 {
		db = db.Select(cfg.selects)
	}

	if len(cfg.omits) > 0 {
		db = db.Omit(cfg.omits...)
	}

	for _, rel := range cfg.preloads {
		db = db.Preload(rel)
	}

	return db
}

var IDOnly = []QueryOption{WithSelect("id")}
