package models

import "errors"

type AliasStatus struct {
	DocCount  int
	Size   int
	FailedIndexOperations int
	RolloverAttemptFailed bool
	Index  string
	Name   string
	Getter AliasGetter
}

type AliasStatuses map[string]AliasStatus


// Refresh() calls GetAlias() and updates the AliasStatus object's DocCount and FailedIndexOperations values 
// to the latest stats for that particular index.
// Returns:
//      - (if err calling .GetAlias()) err from .GetAlias()
//      - nil
func (a *AliasStatus) Refresh() error {
	latestStatus, err := a.Getter.GetAlias(a.Index, a.Name)
	if err != nil {
		return err
	}
	a.DocCount = latestStatus.DocCount
    a.FailedIndexOperations = latestStatus.FailedIndexOperations
	return nil
}

type Difference struct {
    Docs int
    FailedIndexOperations int
}


// GetDifference compares two AliasStatus objects representing the current and previous alias snapshots. 
// If the write index has moved between snapshots it updates the old snapshot to the final values for that index
// and includes any increases between and the original previous snapshot and the updated previous snapshot in the 
// returned Difference.
// Returns:
//    - (if names don't match) err("Cannot compare two different aliases")
//    - (if Refresh methods fails) err from .Refresh()
//    - a Difference object representing the docs added and any new failed index operations between the two snapshots.
func (a AliasStatus) GetDifference(old AliasStatus) (Difference, error) {
    if a.Name != old.Name {
        return Difference{}, errors.New("Cannot compare two different aliases")
    }

    if a.Index != old.Index {
        oldDocCount := old.DocCount
        oldIndexFailureCount := old.FailedIndexOperations
        err := old.Refresh()
        if err != nil {
            return Difference{}, err
        }

        return Difference{
            Docs:             (old.DocCount - oldDocCount) + a.DocCount,
            FailedIndexOperations: (old.FailedIndexOperations - oldIndexFailureCount) + a.FailedIndexOperations,
        }, nil
    }

    return Difference{
        Docs:             a.DocCount - old.DocCount,
        FailedIndexOperations: a.FailedIndexOperations - old.FailedIndexOperations,
    }, nil
}
