package api

import (
	"net/http"
	"time"
)

func v1DeletedResponse() any {
	return struct {
		DeletedAt time.Time `json:"deleted_at"`
		Message   string    `json:"message"`
		Status    int       `json:"status"`
		Version   string    `json:"version"`
	}{
		Version:   "v1",
		DeletedAt: time.Now(),
		Message:   "resource deleted",
		Status:    http.StatusOK,
	}
}

func v1CreatedResponse() any {
	return struct {
		Version   string    `json:"version"`
		CreatedAt time.Time `json:"created_at"`
		Message   string    `json:"message"`
		Status    int       `json:"status"`
	}{
		CreatedAt: time.Now(),
		Message:   "resource created",
		Version:   "v1",
		Status:    http.StatusCreated,
	}
}

func v1NotFoundResponse() any {
	return struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "resource not found",
		Status:  http.StatusNotFound,
	}
}

func v1BadRequestResponse(err error) any {
	return struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "bad request",
		Error:   err.Error(),
		Status:  http.StatusBadRequest,
	}
}

func v1UnprocessableEntityResponse(err error) any {
	return struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "unprocessable entity",
		Error:   err.Error(),
		Status:  http.StatusUnprocessableEntity,
	}
}

func v1InternalServerErrorResponse(err error) any {
	return struct {
		Version string `json:"version"`
		Message string `json:"message"`
		Error   string `json:"error"`
		Status  int    `json:"status"`
	}{
		Version: "v1",
		Message: "internal server error",
		Error:   err.Error(),
		Status:  http.StatusInternalServerError,
	}
}
