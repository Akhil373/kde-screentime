package store

import "testing"

func TestNewStore_CreateTables(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if store.Db == nil {
		t.Fatal("db is nil")
	}

	if store.LastID != nil {
		t.Fatal("LastID is not nil")
	}
}

func TestRecord_InsertsAppAndActivity(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	err = store.Record(WindowInfo{PID: 123, Caption: "test", WMClass: "test", VirtualDesktop: []int32{1}})
	if err != nil {
		t.Fatal(err)
	}

	if store.LastID == nil {
		t.Fatal("LastID is nil")
	}
}

func TestHeartbeat_UpdatesLastID(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	err = store.Record(WindowInfo{PID: 123, Caption: "test", WMClass: "test", VirtualDesktop: []int32{1}})
	if err != nil {
		t.Fatal(err)
	}

	err = store.Heartbeat()
	if err != nil {
		t.Fatal(err)
	}

	if store.LastID == nil {
		t.Fatal("LastID is nil")
	}
}

func TestFinalize_SetsLastID(t *testing.T) {
	store, err := NewStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	err = store.Finalize()
	if err != nil {
		t.Fatal(err)
	}

	if store.LastID != nil {
		t.Fatal("LastID is not nil")
	}
}
