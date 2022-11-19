package notes

import (
	"context"

	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/zerolog/log"
	"github.com/shanmukhsista/go-graphql-starter/pkg/common/db"
	"github.com/shanmukhsista/go-graphql-starter/pkg/common/lib/apperrors"
	"github.com/shanmukhsista/go-graphql-starter/pkg/model"
)

const (
	// use error keys and values to convey errors.use the value for this key within the
	// error_messages.json file to fill in a translation.
	errorFetchingAllNotes          = "errorFetchingAllNotes"
	errorUnknownWhileSavingNewNote = "errorUnknownWhileSavingNewNote"
)

type Service interface {
	SaveNewNote(ctx context.Context, newNoteInput model.NewNoteInput) (*model.Note, error)
	GetAllNotes(ctx context.Context) ([]*model.Note, error)
}

type databaseBackedService struct {
	transactionManager db.TransactionManager
	notesRepository    Repository
}

func (d databaseBackedService) SaveNewNote(ctx context.Context, newNoteInput model.NewNoteInput) (*model.Note, error) {
	// TODO implement me
	noteToSave := newNoteFromRequest(newNoteInput)
	// save this and return
	err := d.transactionManager.WithinTransaction(ctx, func(txnContext context.Context) error {
		savedNote, err := d.notesRepository.CreateNewNote(txnContext, noteToSave)
		if err != nil {
			return err
		}
		log.Debug().Msgf("Saved note with id %s ", savedNote.ID)
		return nil
	})
	if err != nil {
		return nil, apperrors.NewInternalErrorWithUnderlying(errorUnknownWhileSavingNewNote, err)
	}
	return noteToSave, nil
}

func (d databaseBackedService) GetAllNotes(ctx context.Context) ([]*model.Note, error) {
	notes, err := d.notesRepository.GetAllNotes(ctx)
	if err != nil {
		log.Error().Err(err)
		return nil, apperrors.NewInternalErrorWithUnderlying(errorFetchingAllNotes, err)
	}
	return notes, nil
}

// region helper methods

func newNoteFromRequest(req model.NewNoteInput) *model.Note {
	return &model.Note{
		ID:      shortuuid.New(),
		Title:   req.Title,
		Content: req.Content,
	}
}

// endregion

func ProvideNewNotesService(notesRepository Repository, transactionManager db.TransactionManager) (Service, error) {
	return &databaseBackedService{notesRepository: notesRepository,
		transactionManager: transactionManager}, nil
}
