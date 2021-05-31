package blogaggregatormodule

import "github.com/wepala/weos"

type BlogRepository interface {
	weos.Repository
	GetByID(id string) (*Blog)
}