package utils

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func StringToUUID(r *http.Request, param string) (parsedUUID uuid.UUID, err error) {
	paramStr := chi.URLParam(r, param)
	parsedUUID, err = uuid.Parse(paramStr)
	if err != nil {
		return uuid.Nil, err
	}
	return parsedUUID, nil
}
