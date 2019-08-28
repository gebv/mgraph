package mgraph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func NewMGraphPostgres(db *sql.DB) *MGraphPostgres {
	return &MGraphPostgres{db: db}
}

type MGraphPostgres struct {
	db *sql.DB
}

func (s *MGraphPostgres) CreateGraph(ctx context.Context, extID string) (grapID, rootID int64, err error) {
	hashedExtID := md5Hash(strings.ToLower(strings.TrimSpace(extID)))

	tx, err := s.db.Begin()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed open transaction")
	}
	defer tx.Rollback()

	var newGraphID, rootNodeID int64
	err = tx.QueryRow(`INSERT INTO mgraph.graph(hashed_ext_id) VALUES($1) RETURNING graph_id`, hashedExtID).Scan(&newGraphID)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed create new graph")
	}
	err = tx.QueryRow(`INSERT INTO mgraph.graph_nodes(path) VALUES(null) RETURNING node_id`).Scan(&rootNodeID)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed create root node to new graph")
	}
	_, err = tx.Exec(`UPDATE mgraph.graph SET root_node_id = $1 WHERE graph_id = $2`, rootNodeID, newGraphID)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed update graph - set root node ID")
	}
	if err := tx.Commit(); err != nil {
		return 0, 0, errors.Wrap(err, "failed commit transaction")
	}
	return newGraphID, rootNodeID, nil
}

func (s *MGraphPostgres) Add(ctx context.Context, parentID int64) (nodeID int64, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "failed open transaction")
	}
	defer tx.Rollback()
	var parentPath *string
	err = tx.QueryRow(`SELECT path FROM mgraph.graph_nodes WHERE node_id = $1`, parentID).Scan(&parentPath)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrNodeNotFound
		}
		return 0, errors.Wrap(err, "failed find parent node")
	}

	var pathNewNode string
	if parentPath == nil {
		pathNewNode = fmt.Sprint(parentID)
	} else {
		pathNewNode = fmt.Sprint(*parentPath, ".", parentID)
	}

	var newNodeID int64
	if err := tx.QueryRow(`INSERT INTO mgraph.graph_nodes(path, parent_id) VALUES($1, $2) RETURNING node_id`, pathNewNode, parentID).Scan(&newNodeID); err != nil {
		return 0, errors.Wrap(err, "failed insert new node")
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.Wrap(err, "failed commit transaction")
	}
	return newNodeID, nil
}

func (s *MGraphPostgres) Move(ctx context.Context, nodeID, newParentID int64) error {
	_, err := s.db.Exec(`SELECT mgraph.move_node($1, $2)`, nodeID, newParentID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if strings.HasPrefix(err.Where, "PL/pgSQL function mgraph.move_node(bigint,bigint)") &&
				strings.HasPrefix(err.Message, "Not allowed to move in own subthread") {
				return ErrNotAllowedMoveInOwnSubthred
			}

			if strings.HasPrefix(err.Where, "PL/pgSQL function mgraph.move_node(bigint,bigint)") &&
				strings.HasPrefix(err.Message, "Not found parent node") {
				return ErrNodeNotFound
			}
		}
		return err
	}
	return nil
}
func (s *MGraphPostgres) Remove(ctx context.Context, nodeID int64) error {
	_, err := s.db.Exec(`SELECT mgraph.remove_node($1)`, nodeID)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if strings.HasPrefix(err.Where, "PL/pgSQL function mgraph.remove_node(bigint)") &&
				strings.HasPrefix(err.Message, "Not allowed to move in own subthread") {
				return ErrNotAllowedMoveInOwnSubthred
			}

			if strings.HasPrefix(err.Where, "PL/pgSQL function mgraph.remove_node(bigint)") &&
				strings.HasPrefix(err.Message, "Not found parent node") {
				return ErrNodeNotFound
			}
		}
		return err
	}
	return nil
}
func (s *MGraphPostgres) GraphNodes(ctx context.Context, graphID int64) []Node {
	var rootID int64
	if err := s.db.QueryRow(`SELECT root_node_id FROM mgraph.graph WHERE graph_id = $1`, graphID).Scan(&rootID); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("not fodun graph by ID=%d\n", graphID)
			return nil
		}
		log.Println("failed find graph by ID", err)
		return nil
	}
	query := `SELECT node_id, parent_id, path FROM mgraph.graph_nodes WHERE $1::ltree @> path;`
	rows, err := s.db.Query(query, fmt.Sprint(rootID))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("graph=%d is empty\n", graphID)
			return nil
		}
		log.Println("failed find nodes by graph", err)
		return nil
	}
	defer rows.Close()

	var nodeID int64
	var parentID *int64
	var path *string
	var res = []Node{}
	for rows.Next() {
		if err := rows.Scan(&nodeID, &parentID, &path); err != nil {
			log.Println("failed scan row", err)
			return nil
		}
		res = append(res, Node{
			NodeID:   nodeID,
			Path:     path,
			ParentID: parentID,
			Level:    len(strings.Split(strOr(path), ".")),
		})
	}
	return res
}

func strOr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func NewPostgtresConnect(connString string) *sql.DB {
	// Postgres init
	sqlDB, err := sql.Open("postgres", connString)
	if err != nil {
		log.Panicln("failed to connect to PostgreSQL", err)
	}
	sqlDB.SetConnMaxLifetime(0)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(5)
	if err = sqlDB.Ping(); err != nil {
		log.Panicln("failed to connect ping PostgreSQL", err)
	}
	return sqlDB
}
