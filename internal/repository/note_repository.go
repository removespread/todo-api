package repository

import (
	"context"
	"fmt"
	"time"
	"todo-api/internal/models"
	"todo-api/pkg/logger/cache"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	noteCacheKeyPrefix = "note:"
	noteCacheKey       = "notes:all"
	cacheDuration      = 5 * time.Minute
)

type NoteRepository struct {
	pool   *pgxpool.Pool
	cache  cache.Cache
	logger *zap.Logger
}

func NewNoteRepository(pool *pgxpool.Pool, cache cache.Cache, logger *zap.Logger) *NoteRepository {
	return &NoteRepository{
		pool:   pool,
		cache:  cache,
		logger: logger,
	}
}

func (r *NoteRepository) CreateNote(ctx context.Context, note *models.Note) error {
	query := `INSERT INTO notes (title, content) VALIES ($1, $2) RETURNING id`
	err := r.pool.QueryRow(ctx, query, note.Title, note.Content).Scan(&note.ID)
	if err != nil {
		return fmt.Errorf("Failed to create note %w", err)
	}

	if err := r.cache.Delete(ctx, noteCacheKey); err != nil {
		r.logger.Warn("Failed to invalidate note cache", zap.Error(err))
	}

	return nil
}

func (r *NoteRepository) GetAllNotes(ctx context.Context) ([]*models.Note, error) {
	var notes []*models.Note

	// trying get notes from cache
	err := r.cache.Get(ctx, noteCacheKey, &notes)
	if err == nil {
		r.logger.Debug("Notes retrieved from cache")
		return notes, nil
	}

	// if no notes in cache, trying from database
	query := `SELECT id, title, content FROM notes`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to query notes %v", err)
	}

	for rows.Next() {
		note := &models.Note{}
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			return nil, fmt.Errorf("Failed to scan notes: %w", err)
		}
		notes = append(notes, note)
	}

	// saving in cache
	if err := r.cache.Set(ctx, noteCacheKey, notes, cacheDuration); err != nil {
		r.logger.Warn("Failed to cache notes", zap.Error(err))
	}

	return notes, nil
}

func (r *NoteRepository) GetNoteByID(ctx context.Context, id int) (*models.Note, error) {
	note := &models.Note{}
	cacheKey := fmt.Sprintf("%s%d", noteCacheKeyPrefix, id)

	// trying get note from cache
	err := r.cache.Get(ctx, cacheKey, note)
	if err == nil {
		r.logger.Debug("Notes retrieved from cache")
		return note, nil
	}

	// if no notes in cache, trying from database
	query := `SELECT id, title, content FROM notes WHERE id = $1`
	err = r.pool.QueryRow(ctx, query, id).Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		return nil, fmt.Errorf("Failed to get note: %w", err)
	}

	// saving in cache
	if err := r.cache.Set(ctx, cacheKey, note, cacheDuration); err != nil {
		r.logger.Warn("Failed to cache notes", zap.Error(err))
	}

	return note, nil
}

func (r *NoteRepository) UpdateNote(ctx context.Context, note *models.Note) error {
	query := `UPDATE note SET title = $1, content $2 WHERE id = $3`
	result, err := r.pool.Exec(ctx, query, note.Title, note.Content, note.ID)
	if err != nil {
		return fmt.Errorf("Failed to update note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("Note not found")
	}

	// validate cache
	cacheKey := fmt.Sprintf("%s%d", noteCacheKeyPrefix, note.ID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("Failed to invalidate note cache", zap.Error(err))
	}
	if err := r.cache.Delete(ctx, noteCacheKey); err != nil {
		r.logger.Warn("Failed to invalidate notes cache", zap.Error(err))
	}

	return nil
}

func (r *NoteRepository) DeleteNote(ctx context.Context, id int) error {
	query := `DELETE FROM notes WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Failed to delete note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("Note not found")
	}

	cacheKey := fmt.Sprintf("%s%d", noteCacheKeyPrefix, id)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		r.logger.Warn("failed to invalidate note cache", zap.Error(err))
	}
	if err := r.cache.Delete(ctx, noteCacheKey); err != nil {
		r.logger.Warn("failed to invalidate notes cache", zap.Error(err))
	}

	return nil
}
