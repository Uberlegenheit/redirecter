package dao

func (d *DAO) Health() (err error) {

	err = d.db.DB().Ping()
	if err != nil {
		return err
	}

	return nil
}
