package service

import (
	"context"
	"todo-api/internal/models"

	"go.uber.org/zap"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, note *models.Note) error
	GetAllNotes(ctx context.Context) ([]*models.Note, error)
	GetNoteByID(ctx context.Context, id int) (*models.Note, error)
	UpdateNote(ctx context.Context, note *models.Note) error
	DeleteNote(ctx context.Context, id int) error
}

type NoteService struct {
	repo   NoteRepository
	logger *zap.Logger
}

func NewNoteService(repo NoteRepository, logger *zap.Logger) *NoteService {
	return &NoteService{
		repo:   repo,
		logger: logger,
	}
}

func (n *NoteService) CreateNote(ctx context.Context, title string) error {
	n.logger.Info("Creating new note:", zap.String("title", title))

	note := &models.Note{Title: title}

	err := n.repo.CreateNote(ctx, note)
	if err != nil {
		n.logger.Error("Failed to create note", zap.Error(err))
		return err
	}

	n.logger.Info("Note created successfully")
	return nil
}

func (n *NoteService) GetAllNotes(ctx context.Context) ([]*models.Note, error) {
	n.logger.Info("Getting all notes:")

	notes, err := n.repo.GetAllNotes(ctx)
	if err != nil {
		n.logger.Error("Failed to get all notes", zap.Error(err))
		return nil, err
	}
	n.logger.Info("Notes retrieved successfully", zap.Any("Notes", notes))
	return notes, nil
}

func (n *NoteService) GetNoteByID(ctx context.Context, id int) (*models.Note, error) {
	n.logger.Info("Getting note by id:", zap.Int("id", id))

	note, err := n.repo.GetNoteByID(ctx, id)
	if err != nil {
		n.logger.Error("Failed to get note by ID", zap.Error(err))
		return nil, err
	}
	n.logger.Info("Note retrieved successfully", zap.Any("Note", note))
	return note, nil
}

func (n *NoteService) UpdateNote(ctx context.Context, note *models.Note) error {
	n.logger.Info("Updating note", zap.Int("ID", note.ID), zap.String("Title", note.Title))

	existingNote, err := n.repo.GetNoteByID(ctx, note.ID)
	if err != nil {
		n.logger.Error("Note not found", zap.Error(err))
		return err
	}

	existingNote.Title = note.Title
	existingNote.Content = note.Content

	err = n.repo.UpdateNote(ctx, existingNote)
	if err != nil {
		n.logger.Error("Failed to update note", zap.Error(err))
		return err
	}

	n.logger.Info("Note updated successfully", zap.Int("ID", note.ID))
	return nil
}

func (n *NoteService) DeleteNote(ctx context.Context, id int) error {
	n.logger.Info("Deleting note", zap.Int("ID", id))

	_, err := n.repo.GetNoteByID(ctx, id)
	if err != nil {
		n.logger.Error("Note not found", zap.Error(err))
		return err
	}
	err = n.repo.DeleteNote(ctx, id)
	if err != nil {
		n.logger.Error("Failed to delete note", zap.Error(err))
	}

	n.logger.Info("Note deleted successfully", zap.Int("ID", id))
	return nil
}
