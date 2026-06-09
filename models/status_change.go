package models

type StatusChange struct {
	Alias                     string
	DocsAdded                 int
	NewIndexOperationFailures int
}
