package services

func (s *Service) Health() (err error) {

	err = s.dao.Health()
	if err != nil {
		return err
	}

	return nil
}
