package types

import "errors"

type ApiError struct {
	Message string
	Status  int
}

func (e *ApiError) Error() string {
	return e.Message
}

type OtpCodeRequest struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

func (c *OtpCodeRequest) Validate() error {
	if c.Username == "" || c.Code == "" {
		return errors.New("Incorrectly entered data")
	}

	// other checks
	return nil
}
