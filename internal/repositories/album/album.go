package album

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type AlbumRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) (*AlbumRepository, error) {
	return &AlbumRepository{DB: db}, nil
}

func (r *AlbumRepository) Add(ctx context.Context, newAlbum Album) (Album, error) {
	insertRow := r.DB.QueryRowContext(ctx, "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", newAlbum.Title, newAlbum.Artist, newAlbum.Price)

	var id int64
	if err := insertRow.Scan(&id); err != nil {
		return Album{}, fmt.Errorf("Album :: Add :: Insertion: %v", err)
	}

	return r.GetByID(ctx, id)
}

func (r *AlbumRepository) GetByID(ctx context.Context, id int64) (Album, error) {
	var alb Album

	row := r.DB.QueryRowContext(ctx, "SELECT * FROM album WHERE id = $1", id)

	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("Album :: GetByID :: no album with id foung: %d", id)
		}

		return alb, fmt.Errorf("Album :: GetByID :: error for id %d: %v", id, err)
	}

	return alb, nil
}

func (r *AlbumRepository) List(ctx context.Context) ([]Album, error) {
	var albums []Album

	rows, err := r.DB.QueryContext(ctx, "SELECT * FROM album")
	if err != nil {
		zap.L().Error("Album Repository :: Querying all albums from the database failed", zap.Error(err))
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("Album :: List :: Scan Rows: %v", err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Album :: List :: Rows: %v", err)
	}

	return albums, nil
}

func (r *AlbumRepository) ListByArtist(ctx context.Context, name string) ([]Album, error) {
	var albums []Album

	rows, err := r.DB.QueryContext(ctx, "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByPrice %q: %v", name, err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}
