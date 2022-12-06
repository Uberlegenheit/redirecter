package dao

import (
	"everstake-affiliate/models"
)

type DbDAO interface {
	Health() error

	GetLinkByCode(linkCode string) (link models.Link, isFound bool, err error)
}
