package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/NikitaKurabtsev/url-shortener/internal/lib/api/response"
	"github.com/NikitaKurabtsev/url-shortener/internal/lib/logger/sl"
	"github.com/NikitaKurabtsev/url-shortener/internal/lib/random"
	"github.com/NikitaKurabtsev/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.urls.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err = validator.New().Struct(req)
		if err != nil {
			validationErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validationErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, req.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlExists) {
				log.Info("url already exists", sl.Err(err))
				render.JSON(w, r, response.Error("url already exists"))
				return
			}
			log.Error("failed to save", sl.Err(err))
			render.JSON(w, r, response.Error("failed to save url"))
			return
		}

		log.Info("url added", "id", id)

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})

	}
}
