package api

import (
	"encoding/json"
	"net/http"

	"github.com/phillipmugisa/go_user_api/data"
	"github.com/phillipmugisa/go_user_api/types"
)

// util function to write json responses
func writeJsonResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiHandler func(w http.ResponseWriter, r *http.Request) *types.ApiError

// converts custom handler to http.HandlerFunc
func makeHttpHandler(h ApiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			// handle returned error
			writeJsonResponse(w, err.Status, struct {
				Error  string
				Status int
			}{Error: err.Error(), Status: err.Status})
		}
	}
}

// handle user account creation
func (s *ApiServer) handleCreateUser(w http.ResponseWriter, r *http.Request) *types.ApiError {

	user := new(data.User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		return &types.ApiError{
			Message: "Incorrectly entered data",
			Status:  http.StatusNotFound,
		}
	}

	// perform all backend data validations
	v_err := user.Validate()
	if v_err != nil {
		return &types.ApiError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	// check if user exists with submit data
	results, s_err := s.store.GetUser(user.UserName, user.Email)
	if s_err != nil {
		return &types.ApiError{
			Message: "Internal Server Error(an error for all unforeseen situations)",
			Status:  http.StatusInternalServerError,
		}
	}

	if len(results) > 0 {
		return &types.ApiError{
			Message: "A user with this username already exists",
			Status:  http.StatusConflict,
		}
	}

	// hash password
	pwd, err := HashPassword(user.Password)
	if err != nil {
		// password hashing failed
		return &types.ApiError{
			Message: "Internal Server Error(an error for all unforeseen situations)",
			Status:  http.StatusInternalServerError,
		}
	}
	user.Password = pwd

	// save user to db
	results, s_err = s.store.CreateUser(user)
	if s_err != nil {
		return &types.ApiError{
			Message: "Internal Server Error(an error for all unforeseen situations)",
			Status:  http.StatusInternalServerError,
		}
	}

	resp := struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}{
		Username: results[0].UserName,
		Code:     results[0].Code,
	}

	if err := SendEmail(results[0].UserName, results[0].Email, results[0].Code); err != nil {
		// delete user if email was not send
		s.store.DeleteUser(user)
		return &types.ApiError{
			Message: "Internal Server Error(an error for all unforeseen situations)",
			Status:  http.StatusInternalServerError,
		}
	}

	writeJsonResponse(w, http.StatusOK, resp)

	return nil
}

func (s *ApiServer) handleUserVerification(w http.ResponseWriter, r *http.Request) *types.ApiError {
	req_data := new(types.OtpCodeRequest)
	err := json.NewDecoder(r.Body).Decode(req_data)
	if err != nil {
		return &types.ApiError{
			Message: "Incorrectly entered data",
			Status:  http.StatusNotFound,
		}
	}

	if v_err := req_data.Validate(); v_err != nil {
		return &types.ApiError{
			Message: "Incorrectly entered data",
			Status:  http.StatusNotFound,
		}
	}

	// get user with passed user and compare codes
	results, fetch_err := s.store.GetUser(req_data.Username)
	if fetch_err != nil {
		return &types.ApiError{
			Message: "Internal Server Error(an error for all unforeseen situations)",
			Status:  http.StatusInternalServerError,
		}
	}

	// no user with such username was found
	if len(results) == 0 {
		return &types.ApiError{
			Message: "Incorrectly entered data",
			Status:  http.StatusNotFound,
		}
	}

	// compare codes
	if req_data.Code != results[0].Code {
		return &types.ApiError{
			Message: "Incorrectly entered data",
			Status:  http.StatusNotFound,
		}
	}

	// user was found with matching username and code
	// update user details
	results, err = s.store.CompleteUserCheck(results[0].UserName)
	if err != nil {
		return &types.ApiError{
			Message: "Internal Server Error(an error for all unforeseen situations)",
			Status:  http.StatusInternalServerError,
		}
	}

	writeJsonResponse(w, http.StatusOK, struct{}{})

	return nil
}
