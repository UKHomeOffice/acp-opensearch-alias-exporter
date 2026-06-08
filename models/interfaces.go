package models

type Updater interface {
	UpdateHealth(aliasStatuses AliasStatuses)
	UpdateRolloverAttemptFailures(aliasStatuses AliasStatuses)
	UpdateDocsAddedRate([]StatusChange)
	UpdateIndexOperationFailures(statusChanges []StatusChange)
}

type AliasGetter interface {
	GetAlias(index string, name string) (AliasStatus, error)
}
