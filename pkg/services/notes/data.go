package notes

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/shanmukhsista/go-graphql-starter/pkg/common/db"
	"github.com/shanmukhsista/go-graphql-starter/pkg/model"
)

type Repository interface {
	CreateNewNote(ctx context.Context, note *model.Note) (*model.Note, error)
	GetAllNotes(ctx context.Context) ([]*model.Note, error)
	UpdateNote(ctx context.Context, noteId string, note *model.Note) (*model.Note, error)
	ExistsNoteWithID(ctx context.Context, noteId string) (bool, error)
}

// local database models.
type PostgresNotesRepository struct {
	db *db.Database
}

func (p PostgresNotesRepository) CreateNewNote(ctx context.Context, note *model.Note) (*model.Note, error) {
	// TODO implement me
	createNewNoteQuery := `
		insert into notes ( id , title, content ) values ( $1  , $2 , $3)
	`
	tx, err := p.db.GetExistingOrNewTransaction(ctx)
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec(ctx, createNewNoteQuery, note.ID, note.Title, note.Content)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("Inserted a new note . %d rows affected ", res.RowsAffected())
	return note, nil
}

func (p PostgresNotesRepository) GetAllNotes(ctx context.Context) ([]*model.Note, error) {
	selectQuery := `
		select * from notes ; 
	`
	tx, err := p.db.GetExistingOrNewTransaction(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, selectQuery)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	var allNotes = make([]*model.Note, 0)
	for rows.Next() {
		var note model.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			return nil, err
		}
		allNotes = append(allNotes, &note)
	}
	return allNotes, nil
}

func (p PostgresNotesRepository) UpdateNote(ctx context.Context, noteId string, note *model.Note) (*model.Note, error) {
	updateNoteQuery := `
		update notes set title = $1, content = $2  where id = $3
	`
	tx, err := p.db.GetExistingOrNewTransaction(ctx)
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec(ctx, updateNoteQuery, note.Title, note.Content, note.ID)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("Updated notes table . %d rows affected ", res.RowsAffected())
	return note, nil
}

func (p PostgresNotesRepository) ExistsNoteWithID(ctx context.Context, noteId string) (bool, error) {
	var exists bool
	existsQuery := fmt.Sprintf("SELECT exists (%s)", "select id from notes where id = $1")
	err := p.db.QueryRow(ctx, existsQuery, noteId).Scan(&exists)
	if err != nil {
		return exists, err
	}
	return exists, nil
}

func ProvideNewNotesRepository(db *db.Database) (Repository, error) {
	return &PostgresNotesRepository{db: db}, nil
}
