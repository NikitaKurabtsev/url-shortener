package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/NikitaKurabtsev/url-shortener/internal/lib/api/response"
	"github.com/NikitaKurabtsev/url-shortener/internal/lib/logger/sl"
	"github.com/NikitaKurabtsev/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlSaver URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.urls.redirect.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("alias is empty"))

			return
		}

		url, err := urlSaver.GetURL(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "alias", alias)

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to fetch url", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("found url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
