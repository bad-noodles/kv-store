package store

import (
	"os"
	"testing"

	"github.com/bad-noodles/kv-store/pkg/store"
	typesystem "github.com/bad-noodles/kv-store/pkg/type_system"
)

func mustStore(t *testing.T) *store.Store {
	f, err := os.CreateTemp("", "wal")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())

	return store.NewStore(f.Name())
}

func fail(t *testing.T, resp typesystem.Type) {
	t.Fatal("failed with ", resp.Pretty())
}

func TestSetAndGet(t *testing.T) {
	store := mustStore(t)

	resp := store.ExecuteCommand("SET x \"test\"")

	if resp.Value() != "+OK" {
		fail(t, resp)
	}

	resp = store.ExecuteCommand("GET x")

	if resp.Value() != "test" {
		fail(t, resp)
	}
}

func TestGetArray(t *testing.T) {
	store := mustStore(t)

	resp := store.ExecuteCommand("SET x [\"test\"]")

	if resp.Value() != "+OK" {
		fail(t, resp)
	}

	resp = store.ExecuteCommand("GET x")

	switch arr := resp.Value().(type) {
	case []typesystem.Type:
		if arr[0].Value() != "test" {
			t.Fail()
		}
	default:
		t.Fatal("Expected an array")
	}
}

func TestGetInexistent(t *testing.T) {
	store := mustStore(t)

	resp := store.ExecuteCommand("GET x")

	if resp.Value() != nil {
		fail(t, resp)
	}
}
