package retriever

import (
	"testing"
)

func TestIndexgetter_GetIndexes_no_write_indexes(t *testing.T) {
	getter := func(string, string, string) ([]byte, error) {
		return []byte(`
{"kubernetes-acp-prod-egar-production-000102": {
    "aliases": {
      "kubernetes-acp-prod-egar-production": {
        "is_write_index": false
      }
    }
  },
  ".kibana_96354_abc_2": {
    "aliases": {}
  }}`), nil
	}
	sut := NewIndexGetter("foo", "bar", "abc", getter)
	indexes, err := sut.GetWriteIndexes()
	if err != nil {
		t.Error("Expected nil got error")
	}
	if len(indexes) != 0 {
		t.Error("Expected 0 aliases got ", len(indexes))
	}
}

func TestIndexgetter_GetIndexes_one_write_index(t *testing.T) {
	getter := func(string, string, string) ([]byte, error) {
		return []byte(`
{"kubernetes-acp-prod-egar-production-000102": {
    "aliases": {
      "kubernetes-acp-prod-egar-production": {
        "is_write_index": true
      }
    }
  },
  ".kibana_96354_abc_2": {
    "aliases": {}
  }}`), nil
	}
	sut := NewIndexGetter("foo", "bar", "abc", getter)
	indexes, err := sut.GetWriteIndexes()
	if err != nil {
		t.Error("Expected nil got error")
	}
	if len(indexes) != 1 {
		t.Error("Expected 1 aliases got ", len(indexes))
	}
}
