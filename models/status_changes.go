package models


// GetStatusChanges compares two slices of AliasStatus objects representing the (old) previous and (new) current alias status snapshots.
// Returns:
// 		- (if no err) a slice of StatusChange objects that describe how many docs have been added 
// 			and how many index operations have failed between snapshots.
// 		- (if err) an error if any individual alias diff fails.
func GetStatusChanges(oldStatuses AliasStatuses, newStatuses AliasStatuses) ([]StatusChange, error) {
	statusChanges := make([]StatusChange, 0)
	for _, newStatus := range newStatuses {
		if oldStatus, ok := oldStatuses[newStatus.Name]; ok {
			difference, err := newStatus.GetDifference(oldStatus)
			if err != nil {
				return nil, err
			}
			statusChanges = append(statusChanges, StatusChange{Alias: newStatus.Name, DocsAdded: difference.Docs, NewIndexOperationFailures: difference.FailedIndexOperations})
		} else {
			statusChanges = append(statusChanges, StatusChange{Alias: newStatus.Name, DocsAdded: newStatus.DocCount, NewIndexOperationFailures: newStatus.FailedIndexOperations})
		}
	}
	return statusChanges, nil
}
