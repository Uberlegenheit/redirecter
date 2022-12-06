package dao

import (
	"everstake-affiliate/conf"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DAO struct {
	db *gorm.DB
}

func NewDAO(c conf.Config) (DbDAO, error) {
	db, err := gorm.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
			c.Postgres.User,
			c.Postgres.Password,
			c.Postgres.Host,
			c.Postgres.Database,
			c.Postgres.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	return New(db), nil

}

func New(db *gorm.DB) *DAO {
	return &DAO{
		db: db,
	}
}

type TXDAO struct {
	DAO
}
