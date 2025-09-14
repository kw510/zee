CREATE TABLE tuples (
  parent_namespace TEXT NOT NULL,
  parent_id TEXT NOT NULL,
  parent_relation TEXT NOT NULL,
  child_namespace TEXT NOT NULL,
  child_id TEXT NOT NULL,
  child_relation TEXT NOT NULL,

  PRIMARY KEY (parent_namespace, parent_id, parent_relation, child_namespace, child_id, child_relation)
);

CREATE INDEX idx_tuples_parent ON tuples (parent_namespace, parent_id, parent_relation);
CREATE INDEX idx_tuples_child ON tuples (child_namespace, child_id, child_relation);