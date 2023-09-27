package models

type Updater interface {
	Update([]CountRate)
}

type AliasGetter interface {
	GetAlias(index string, name string) (AliasStatus, error)
}
