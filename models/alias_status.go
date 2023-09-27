package models

import "errors"

type AliasStatus struct {
	Count  int
	Size   int
	Failed bool
	Index  string
	Name   string
	Getter AliasGetter
}

type AliasStatuses map[string]AliasStatus

func (a *AliasStatus) Refresh() error {
	as, err := a.Getter.GetAlias(a.Index, a.Name)
	if err != nil {
		return err
	}
	a.Count = as.Count
	return nil
}

func (a AliasStatus) Diff(old AliasStatus) (int, error) {
	if a.Name != old.Name {
		return 0, errors.New("Cannot compare two different aliases")
	}

	if a.Index != old.Index {
		oldCount := old.Count
		err := old.Refresh()
		if err != nil {
			return 0, err
		}

		return (old.Count - oldCount) + a.Count, nil
	}

	return a.Count - old.Count, nil
}
