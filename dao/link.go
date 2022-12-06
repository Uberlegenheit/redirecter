package dao

import (
	"everstake-affiliate/models"
	"github.com/jinzhu/gorm"
)

func (d *DAO) GetLinkByCode(linkCode string) (link models.Link, isFound bool, err error) {

	err = d.db.Model(&models.Link{}).
		Where("code = ?", linkCode).
		Find(&link).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return link, false, nil
		}

		return link, false, err
	}

	return link, true, nil
}
