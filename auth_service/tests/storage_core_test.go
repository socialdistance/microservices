package tests

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"grpc_auth_tutorial/sso/internal/domain/models"
	"grpc_auth_tutorial/sso/internal/storage/core"
)

func TestPut(t *testing.T) {
	store := core.NewCoreStore()

	key := uuid.New()
	value := models.StoreEvent{
		ID:        key,
		Title:     "test",
		Started:   time.Now(),
		Ended:     time.Now(),
		ServiceID: 5,
	}

	// defer delete(store, key)

	// // Sanity check
	// _, contains = store.Store[key]
	// if contains {
	// 	t.Error("key/value already exists")
	// }

	// err should be nil
	err := store.Put(key, value)
	if err != nil {
		t.Error(err)
	}

	val, err := store.Get(key)
	if err != nil {
		t.Error("create failed")
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}

// func TestGet(t *testing.T) {
// 	const key = "read-key"
// 	const value = "read-value"

// 	var val interface{}
// 	var err error

// 	defer delete(store.m, key)

// 	// Read a non-thing
// 	val, err = Get(key)
// 	if err == nil {
// 		t.Error("expected an error")
// 	}
// 	if !errors.Is(err, ErrorNoSuchKey) {
// 		t.Error("unexpected error:", err)
// 	}

// 	store.m[key] = value

// 	val, err = Get(key)
// 	if err != nil {
// 		t.Error("unexpected error:", err)
// 	}

// 	if val != value {
// 		t.Error("val/value mismatch")
// 	}
// }

// func TestDelete(t *testing.T) {
// 	const key = "delete-key"
// 	const value = "delete-value"

// 	var contains bool

// 	defer delete(store.m, key)

// 	store.m[key] = value

// 	_, contains = store.m[key]
// 	if !contains {
// 		t.Error("key/value doesn't exist")
// 	}

// 	Delete(key)

// 	_, contains = store.m[key]
// 	if contains {
// 		t.Error("Delete failed")
// 	}
// }
