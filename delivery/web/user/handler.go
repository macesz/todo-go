package user

import (
	"encoding/json" // For JSON (like JSON.parse/stringify in JS)
	"errors"
	"net/http" // Standard HTTP library (like fetch in JS or HttpServlet in Java)
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	validate "github.com/go-playground/validator/v10" // For struct validation (like Joi in JS or Hibernate Validator in Java)
	"github.com/macesz/todo-go/delivery/web/utils"
	"github.com/macesz/todo-go/domain"
)

// CreateUser creates a new HTTP handler for creating a new user.
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqUser domain.CreateUserDTO // Empty User struct to decode into

	// Decode the JSON body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		// domain.ErrorResponse{Error: err.Error() for dynamic error message
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	if err := validate.New().Struct(reqUser); err != nil {
		useErr := translateValidationError(err)
		// Dynamic message, e.g., "Name is required; Email is required"
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: useErr})
		return
	}

	// Create the user using the service
	user, err := h.Service.CreateUser(r.Context(), reqUser.Name, reqUser.Email, reqUser.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmail):
			utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
			return
		case errors.Is(err, domain.ErrInvalidPassword):
			utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
			return
		case errors.Is(err, domain.ErrDuplicate):
			utils.WriteJSON(w, http.StatusConflict, domain.ErrorResponse{Error: err.Error()})
			return
		default:
			utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
			return
		}
	}

	respUser := domain.UserResponseDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	utils.WriteJSON(w, http.StatusCreated, respUser)
}

// GetUser creates a new HTTP handler for getting a user by ID.
func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id") // Get the "id" URL parameter

	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	user, err := h.Service.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			utils.WriteJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
		return
	}

	respUser := domain.UserResponseDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	utils.WriteJSON(w, http.StatusOK, respUser)
}

// DeleteUser creates a new HTTP handler for deleting a user.
func (h *UserHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idr := chi.URLParam(r, "id") // Get the "id" URL parameter

	if idr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.ParseInt(idr, 10, 64) // Convert id string to int
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, domain.ErrorResponse{Error: "id must be an integer"})
		return
	}

	err = h.Service.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			utils.WriteJSON(w, http.StatusNotFound, domain.ErrorResponse{Error: err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, domain.ErrorResponse{Error: "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content on successful deletion
}

func (h *UserHandlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Implementation goes here
}

func translateValidationError(err error) string {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return "validation failed"
	}

	messages := []string{}
	for _, fieldErr := range validationErrs {
		switch fieldErr.Field() {
		case "Name":
			switch fieldErr.Tag() {
			case "required":
				messages = append(messages, "Name is required")
			case "min":
				messages = append(messages, "Name must be at least 5 characters")
			case "max":
				messages = append(messages, "Name must be at most 255 characters")
			}
		case "Email":
			switch fieldErr.Tag() {
			case "required":
				messages = append(messages, "Email is required")
			case "min":
				messages = append(messages, "Email must be at least 5 characters")
			case "max":
				messages = append(messages, "Email must be at most 255 characters")
			}
		case "Password":
			switch fieldErr.Tag() {
			case "required":
				messages = append(messages, "Password is required")
			case "min":
				messages = append(messages, "Password must be at least 5 characters")
			case "max":
				messages = append(messages, "Password must be at most 255 characters")
			}
		default:
			messages = append(messages, fieldErr.Field()+" is invalid")
		}
		if len(messages) == 0 {
			return "validation failed"
		}
	}
	return strings.Join(messages, "; ") // Combine if multiple errors

}
