package application

// Demo добавляет пример студентов, критериев и оценок, не удаляя старые записи
func (s *Service) Demo() error {
	anna, err := s.AddStudent("Анна Иванова")
	if err != nil {
		return err
	}
	boris, err := s.AddStudent("Борис Петров")
	if err != nil {
		return err
	}
	homework, err := s.AddCriterion("Домашняя работа", 10, 0.4)
	if err != nil {
		return err
	}
	exam, err := s.AddCriterion("Экзамен", 20, 0.6)
	if err != nil {
		return err
	}
	grades := []struct {
		studentID   int
		criterionID int
		score       float64
	}{
		{anna.ID(), homework.ID(), 8},
		{anna.ID(), exam.ID(), 18},
		{boris.ID(), homework.ID(), 6},
		{boris.ID(), exam.ID(), 10},
	}
	for _, grade := range grades {
		if err := s.SetEvaluation(grade.studentID, grade.criterionID, grade.score); err != nil {
			return err
		}
	}
	return nil
}
