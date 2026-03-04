package service

type ImportResult struct {
	Created   int
	Updated   int
	Unchanged int
	Skipped   int
}
