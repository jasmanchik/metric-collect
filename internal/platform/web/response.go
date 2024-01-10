package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Response(w http.ResponseWriter, data interface{}, statusCode int) error {
	if data == nil && statusCode == http.StatusOK {
		w.WriteHeader(http.StatusNoContent) //todo http: superfluous response.WriteHeader call
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshalling data: %v, %w", data, err)
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		return fmt.Errorf("writing data: %v, %w", data, err)
	}
	return nil
}

func RespondError(w http.ResponseWriter) error {
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Response(w, er, http.StatusInternalServerError); err != nil {
		return err
	}

	return nil
}
