package models

func GetCountChanges(oldStatuses AliasStatuses, newStatuses AliasStatuses) ([]CountRate, error) {
	countChanges := make([]CountRate, 0)
	for _, newStatus := range newStatuses {
		if oldStatus, ok := oldStatuses[newStatus.Name]; ok {
			countChange, err := newStatus.Diff(oldStatus)
			if err != nil {
				return nil, err
			}
			countChanges = append(countChanges, CountRate{Alias: newStatus.Name, Total: countChange})
		} else {
			countChanges = append(countChanges, CountRate{Alias: newStatus.Name, Total: newStatus.Count})
		}
	}
	return countChanges, nil
}
