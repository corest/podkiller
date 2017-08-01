package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestSliceOperation(t *testing.T) {

	tables := []struct {
		operation string
		sl1       []string
		sl2       []string
		result    []string
	}{
		{operation: "substruction",
			sl1:    []string{"apple", "banana", "perry"},
			sl2:    []string{"apple", "cherry", "strawberry"},
			result: []string{"banana", "perry"},
		},
		{operation: "unity",
			sl1:    []string{"apple", "banana", "perry"},
			sl2:    []string{"apple", "cherry", "strawberry"},
			result: []string{"apple"},
		},
	}

	for _, table := range tables {
		got, _ := sliceOperation(table.operation, table.sl1, table.sl2)
		want := table.result

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Slice operation was incorrect, got: [%s], want: [%s]",
				strings.Join(want, ", "), strings.Join(want, ", "))
		}
	}

	if _, err := sliceOperation("non-supported", tables[0].sl1, tables[0].sl2); err == nil {
		t.Error("Unsupported slice operation should fail")
	}

}
