package apierror

import "fmt"

type ApiError struct {
	Description  string `json:"description"`
	Err          error  `json:"-"`
	DeveloperMsg string `json:"developer_msg"`
	Code         string `json:"code"`
}

func New(description, developerMsg, code string) *ApiError {
	return &ApiError{
		Description:  description,
		Err:          fmt.Errorf(description),
		DeveloperMsg: developerMsg,
		Code:         code,
	}
}
