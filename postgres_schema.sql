BEGIN;

CREATE EXTENSION IF NOT EXISTS ltree;

CREATE SCHEMA mgraph;

CREATE TABLE mgraph.graph_nodes (
    node_id bigserial PRIMARY KEY,
    parent_id bigint REFERENCES mgraph.graph_nodes(node_id),
    path ltree
);

CREATE INDEX graph_nodes_path_gist_idx ON mgraph.graph_nodes USING GIST (path);

COMMENT ON TABLE mgraph.graph_nodes IS 'Таблица с узлами графа (materialized path).';
COMMENT ON COLUMN mgraph.graph_nodes.node_id IS 'Уникальный внутренний идентификатор ноды.';
COMMENT ON COLUMN mgraph.graph_nodes.path IS 'Путь ноды. Первым элементом идентификатор графа.';

CREATE TABLE mgraph.graph (
    graph_id bigserial PRIMARY KEY,
    root_node_id bigint REFERENCES mgraph.graph_nodes(node_id),
    hashed_ext_id varchar NOT NULL,
    CONSTRAINT graph_uniq_ext_id_idx UNIQUE(hashed_ext_id)
);

COMMENT ON TABLE mgraph.graph IS 'Таблица с графами.';
COMMENT ON COLUMN mgraph.graph.root_node_id IS 'Идентификатор корневого узла графа.';
COMMENT ON COLUMN mgraph.graph.hashed_ext_id IS 'Уникальный внешний идентификатор (хеш) графа.';

-- remove node by ID
CREATE OR REPLACE FUNCTION mgraph.remove_node(
    remove_node_id bigint
) RETURNS void  AS $$
    DECLARE
        _path ltree;
    BEGIN
        SELECT gn.path INTO _path FROM mgraph.graph_nodes gn WHERE gn.node_id = remove_node_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Node not found';
        END IF;

        IF _path IS NULL  THEN
            RAISE EXCEPTION 'Not allowed remove root node of graph';
        END IF;

        -- deletes all children
        DELETE FROM mgraph.graph_nodes gn WHERE _path @> gn.path AND _path != gn.path;
        DELETE FROM mgraph.graph_nodes gn WHERE gn.node_id = remove_node_id;
    END;
$$ language plpgsql;

-- move node to new parent node ID
CREATE OR REPLACE FUNCTION mgraph.move_node(
    currnet_node_id bigint,
    new_parent_id bigint
) RETURNS void  AS $$
    DECLARE
        _from_path ltree;
        _to_path ltree;
        _prefix ltree;
    BEGIN
        SELECT gn.path INTO _from_path FROM mgraph.graph_nodes gn WHERE gn.node_id = currnet_node_id;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Node not found current node';
        END IF;

        IF _from_path IS NULL  THEN
            RAISE EXCEPTION 'Not allowed move root node of graph';
        END IF;

        SELECT gn.path INTO _to_path FROM mgraph.graph_nodes gn WHERE gn.node_id = new_parent_id;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Node not found parent node';
        END IF;

        IF _to_path <> NULL AND _from_path::ltree @> _to_path::ltree THEN
            RAISE EXCEPTION 'Not allowed to move in their own subthread';
        END IF;

        -- NOTE: allowed move to root node

        _prefix := coalesce(_to_path, ''::ltree) || new_parent_id::text::ltree;
        UPDATE mgraph.graph_nodes SET parent_id = new_parent_id, path = _prefix WHERE node_id = currnet_node_id;

        -- updates all children and current node
        UPDATE mgraph.graph_nodes gn SET
            path = _prefix ||
                CASE WHEN
                    nlevel(_from_path) >= nlevel(path)
                THEN
                    ''::ltree
                ELSE
                    subpath(path, nlevel(_from_path) , nlevel(path))
                END
        WHERE _from_path || currnet_node_id::text::ltree @> gn.path;
    END;
$$ language plpgsql;

CREATE VIEW v_nodes AS SELECT * FROM graph_nodes ORDER BY nlevel(coalesce(path, ''::ltree)), path;
-- SELECT * FROM v_nodes;

DELETE FROM graph;
DELETE FROM graph_nodes;
SELECT * FROM v_nodes;

COMMIT;
