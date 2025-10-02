package vp

import (
	"context"
	"fmt"

	"github.com/bbars/whispar/internal/tree"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Project struct {
	db *sqlx.DB
}

func Open(file string) (*Project, error) {
	db, err := sqlx.Open("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("error opening project file sqlite3 database: %w", err)
	}

	return &Project{
		db: db,
	}, nil
}

func (p *Project) Close() error {
	return p.db.Close()
}

func (p *Project) InsertModelElement(ctx context.Context, n *tree.Node) error {
	entry, err := newModelElement(ctx, n)
	if err != nil {
		return fmt.Errorf("prepare entry: %w", err)
	}

	_, err = p.db.NamedExecContext(ctx, `
		INSERT INTO MODEL_ELEMENT (
			ID,
			USER_ID,
			USER_ID_PARENT,
			MODEL_TYPE,
			PARENT_ID,
			NAME,
			DEFINITION,
			MIRROR_SOURCE,
			AUTHOR,
			CREATE_AT,
			LAST_MOD_AT
		)
		VALUES (
			:ID,
			:USER_ID,
			:USER_ID_PARENT,
			:MODEL_TYPE,
			:PARENT_ID,
			:NAME,
			:DEFINITION,
			:MIRROR_SOURCE,
			:AUTHOR,
			:CREATE_AT,
			:LAST_MOD_AT
		)
	`, entry)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	return nil
}

func (p *Project) InsertDiagramElement(ctx context.Context, n *tree.Node) error {
	entry, err := newDiagramElement(ctx, n)
	if err != nil {
		return fmt.Errorf("prepare entry: %w", err)
	}

	_, err = p.db.NamedExecContext(ctx, `
		INSERT INTO DIAGRAM_ELEMENT (
			ID,
			SHAPE_TYPE,
			DIAGRAM_ID,
			MODEL_ELEMENT_ID,
			COMPOSITE_MODEL_ELEMENT_ADDRESS,
			REF_MODEL_ELEMENT_ADDRESS,
			PARENT_ID,
			DEFINITION
		)
		VALUES (
			:ID,
			:SHAPE_TYPE,
			:DIAGRAM_ID,
			:MODEL_ELEMENT_ID,
			:COMPOSITE_MODEL_ELEMENT_ADDRESS,
			:REF_MODEL_ELEMENT_ADDRESS,
			:PARENT_ID,
			:DEFINITION
		)
	`, entry)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	return nil
}

func (p *Project) InsertDiagram(ctx context.Context, n *tree.Node) error {
	entry, err := newDiagram(ctx, n)
	if err != nil {
		return fmt.Errorf("prepare entry: %w", err)
	}

	_, err = p.db.NamedExecContext(ctx, `
		INSERT INTO DIAGRAM (
			ID,
			DIAGRAM_TYPE,
			PARENT_MODEL_ID,
			NAME,
			DEFINITION
		)
		VALUES (
			:ID,
			:DIAGRAM_TYPE,
			:PARENT_MODEL_ID,
			:NAME,
			:DEFINITION
		)
	`, entry)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	return nil
}
