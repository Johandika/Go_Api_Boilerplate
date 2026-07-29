package shared

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"strconv"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

type PaginatedResponse struct {
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message"`
	Data         any    `json:"data"`
	TotalData    int64  `json:"total_data"`
	TotalPage    int64  `json:"total_page"`
	ActivePage   int32  `json:"active_page"`
	LimitPerPage int32  `json:"limit_per_page"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func WriteSuccess(w http.ResponseWriter, status int, message string, data any) {
	writeJSON(w, status, Response{StatusCode: status, Message: message, Data: data})
}

func WritePaginated(w http.ResponseWriter, message string, data any, totalData int64, page int32, limit int32) {
	totalPage := int64(math.Ceil(float64(totalData) / float64(limit)))
	writeJSON(w, http.StatusOK, PaginatedResponse{
		StatusCode:   http.StatusOK,
		Message:      message,
		Data:         data,
		TotalData:    totalData,
		TotalPage:    totalPage,
		ActivePage:   page,
		LimitPerPage: limit,
	})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{StatusCode: status, Message: message})
}

func WriteInternalError(w http.ResponseWriter, err error) {
	slog.Error("internal server error", "error", err)
	WriteError(w, http.StatusInternalServerError, "internal server error")
}

func DecodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func Pagination(r *http.Request) (int32, int32, int32) {
	limitStr := r.URL.Query().Get("limit_per_page")
	if limitStr == "" {
		limitStr = r.URL.Query().Get("limit")
	}
	limit := parseInt32(limitStr, 10)
	page := parseInt32(r.URL.Query().Get("page"), 1)
	if limit <= 0 || limit > 1000 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	return limit, (page - 1) * limit, page
}

func Search(r *http.Request) string {
	return r.URL.Query().Get("q")
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func parseInt32(value string, fallback int32) int32 {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return fallback
	}
	return int32(parsed)
}
