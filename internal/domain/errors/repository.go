package customerrors

type RepositoryError struct {
	Code        int
	Message     string
	DebugDetail string
}

func (r *RepositoryError) Error() string {
	return r.Message
}

func NotFoundInRepository() *RepositoryError {
	return &RepositoryError{
		Code:    404,
		Message: "not found",
	}
}

func UnknownErrorInRepository(detail ...string) *RepositoryError {
	if len(detail) > 0 {
		return &RepositoryError{
			Code:        500,
			Message:     "unknown error",
			DebugDetail: detail[0],
		}
	}
	return &RepositoryError{
		Code:    500,
		Message: "unknown error",
	}
}

func (r *RepositoryError) ToDomain() *DomainError {
	return &DomainError{
		Code:        r.Code,
		Message:     r.Message,
		DebugDetail: r.DebugDetail,
	}
}
