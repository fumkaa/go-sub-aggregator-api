package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

var v = validator.New()
var decoder = schema.NewDecoder()

func ValidateQuery(log *slog.Logger, model any) func(http.Handler) http.Handler {
	modelType := reflect.TypeOf(model)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.ValidateQuery()"
			log = log.With(
				slog.String("op", op),
			)

			newModel := reflect.New(modelType).Interface()

			if err := decoder.Decode(newModel, r.URL.Query()); err != nil {
				log.ErrorContext(r.Context(), "query decode failed",
					slog.String("query", r.URL.RawQuery),
					slog.String("error", err.Error()),
				)
				http.Error(w, "Invalid query parameters", http.StatusBadRequest)
				return
			}

			if err := v.Struct(newModel); err != nil {
				log.WarnContext(r.Context(), "query validation failed",
					slog.String("error", err.Error()),
				)
				http.Error(w, "Query validation failed: "+err.Error(), http.StatusUnprocessableEntity)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func BindAndValidate(log *slog.Logger, model any) func(http.Handler) http.Handler {
	modelType := reflect.TypeOf(model)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.BindAndValidate()"
			log = log.With(
				slog.String("op", op),
			)

			newModel := reflect.New(modelType).Interface()

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.ErrorContext(r.Context(), "body read failed", slog.String("error", err.Error()))
				http.Error(w, "Read body error", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if err := json.Unmarshal(bodyBytes, newModel); err != nil {
				log.ErrorContext(r.Context(), "json unmarshal failed", slog.String("error", err.Error()))
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if err := v.Struct(newModel); err != nil {
				log.ErrorContext(r.Context(), "body validation failed", slog.String("error", err.Error()))
				http.Error(w, "Validation failed: "+err.Error(), http.StatusUnprocessableEntity)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ValidateUUID(log *slog.Logger, paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.ValidateUUID()"
			log = log.With(
				slog.String("op", op),
			)

			idStr := chi.URLParam(r, paramName)
			if _, err := uuid.Parse(idStr); err != nil {
				log.ErrorContext(r.Context(), "uuid validation failed",
					slog.String("param", paramName),
					slog.String("value", idStr),
				)
				http.Error(w, "invalid UUID format", http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func init() {
	v.RegisterValidation("mm_yyyy", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		_, err := time.Parse("01-2006", dateStr)
		return err == nil
	})
}
