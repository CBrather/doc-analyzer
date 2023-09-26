package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/CBrather/analyzer/internal/api/middleware"
	"github.com/CBrather/analyzer/internal/config"
	infra "github.com/CBrather/analyzer/internal/infrastructure"
	"github.com/CBrather/analyzer/internal/repositories/album"
)

func SetupAlbumRoutes(rootRouter *chi.Mux) {
	config := config.GetEnvironment()

	albumRouter := chi.NewRouter()
	albumRouter.Use(middleware.EnsureValidToken(config))

	db, err := infra.GetDB()
	if err != nil {
		zap.L().Fatal("Failed getting the database while setting up album routes", zap.Error(err))
	}

	repo, err := album.NewRepository(db)
	if err != nil {
		zap.L().Fatal("Failed creating Album Repository.", zap.Error(err))
	}

	albumRouter.With(middleware.RequireScope("read:albums")).Get("/{id}", getAlbum(repo))
	albumRouter.With(middleware.RequireScope("read:albums")).Get("/", listAlbums(repo))
	albumRouter.With(middleware.RequireScope("create:albums")).Post("/", postAlbum(repo))

	rootRouter.Mount("/albums", albumRouter)
}

func listAlbums(repo *album.AlbumRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		albums, err := repo.List(req.Context())
		if err != nil {
			zap.L().Error("Failed to retrieve the list of albums.", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		zap.L().Info("Successfully retrieved a list of albums to return")

		body, err := json.Marshal(albums)
		if err != nil {
			zap.L().Error("Failed to serialize list of albums", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

func getAlbum(repo *album.AlbumRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := chi.URLParam(req, "id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(w, "Invalid id provided", http.StatusBadRequest)
			return
		}

		id64 := int64(id)
		album, err := repo.GetByID(req.Context(), id64)

		if err != nil {
			http.Error(w, "No album with that id was found", http.StatusNotFound)
			return
		}

		body, err := json.Marshal(album)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Failed to serialize the album with id %s", album.ID), zap.Error(err))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

func postAlbum(repo *album.AlbumRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var newAlbum album.Album

		rawRequestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			zap.L().Warn("Unable to read bytes of the request body", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(rawRequestBody, &newAlbum)
		if err != nil {
			zap.L().Info("Unable to deserialize request body to album", zap.Error(err))

			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		addedAlbum, err := repo.Add(req.Context(), newAlbum.Title)
		if err != nil {
			zap.L().Error("Failed to save new album", zap.Error(err))
			zap.L().Debug("Failed to save new album", zap.Any("struct", newAlbum))

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		serializedAlbum, err := json.Marshal(addedAlbum)
		if err != nil {
			zap.L().Error("Failed serializing new album after successful save", zap.Error(err))
			zap.L().Debug("Failed serializing new album after successful save", zap.Any("struct", addedAlbum))

			http.Error(w, "Internal Server Error occurred after the album was successfully saved", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(serializedAlbum)
	}
}
