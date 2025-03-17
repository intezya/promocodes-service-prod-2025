package types

import "time"

type SolutionDate string

func (s *SolutionDate) ToDate() (time.Time, error) {
	if s == nil {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", string(*s))
}
