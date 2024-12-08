package note

import (
	"context"
	"fmt"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/dto"
)

type NoteService struct {
	bookRepo      domain.BookRepository
	noteRepo      domain.NoteRepository
	validationSvc domain.ValidationService
}

func NewNoteService(bookRepo domain.BookRepository, noteRepo domain.NoteRepository, validationSvc domain.ValidationService) *NoteService {
	return &NoteService{
		bookRepo:      bookRepo,
		noteRepo:      noteRepo,
		validationSvc: validationSvc,
	}
}

func (s *NoteService) GetNotes(ctx context.Context, userID, bookID int64) ([]dto.NoteResponse, error) {
	notes, err := s.noteRepo.GetNotesByUserIDAndBookID(ctx, userID, bookID)
	if err != nil {
		return []dto.NoteResponse{}, err
	}

	res := make([]dto.NoteResponse, len(notes))
	for i, note := range notes {
		res[i] = dto.NoteResponse{
			ID:         note.ID,
			PageNumber: note.PageNumber,
			Content:    note.Content,
		}
	}

	return res, nil
}

func (s *NoteService) CreateNote(ctx context.Context, userID, bookID int64, note dto.NoteRequest) (dto.NoteResponse, error) {
	_, err := s.bookRepo.GetBookByUserID(ctx, userID, bookID)
	if err != nil {
		return dto.NoteResponse{}, err
	}

	newNote := domain.Note{
		UserID:     userID,
		BookID:     bookID,
		PageNumber: note.PageNumber,
		Content:    note.Content,
	}
	if err := s.validationSvc.ValidateStruct(newNote); err != nil {
		return dto.NoteResponse{}, err
	}

	createdNote, err := s.noteRepo.CreateNote(ctx, newNote)
	if err != nil {
		return dto.NoteResponse{}, err
	}

	return dto.NoteResponse{
		ID:         createdNote.ID,
		PageNumber: createdNote.PageNumber,
		Content:    createdNote.Content,
	}, nil
}

func (s *NoteService) DeleteNote(ctx context.Context, userID, noteID int64) error {
	note, err := s.noteRepo.GetNoteByUserID(ctx, noteID, userID)
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrRecordNotFound, "note")
	}

	if err := s.noteRepo.DeleteNote(ctx, note.ID); err != nil {
		return err
	}

	return nil
}
