package service

// CommitReport is the reconciliation printed after a commit: the per-source
// counts, the source paths skipped because they had no store, and any
// source-vs-source collisions that made the commit refuse and write nothing.
type CommitReport struct {
	Sources    []SourceReconcile
	Skipped    []string
	Collisions []string
}

// SourceReconcile is one source's commit count: Imported is how many issues were
// written and Skipped is how many were left untouched.
type SourceReconcile struct {
	Path     string
	Imported int
	Skipped  int
}
