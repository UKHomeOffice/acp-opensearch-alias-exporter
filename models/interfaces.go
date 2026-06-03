package models

type Updater interface {
	UpdateCountRates([]CountRate)
	UpdateRolloverHealth(aliasStatuses AliasStatuses)
	UpdateShardHealth(aliasStatuses AliasStatuses)
	UpdateIndexOperationHealth(aliasStatuses AliasStatuses)
}

type AliasGetter interface {
	GetAlias(index string, name string) (AliasStatus, error)
	
}
