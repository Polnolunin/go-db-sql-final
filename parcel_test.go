package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// get
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Client, parcel.Client)
	assert.Equal(t, res.Status, parcel.Status)
	assert.Equal(t, res.Address, parcel.Address)
	assert.Equal(t, res.CreatedAt, parcel.CreatedAt)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)
	_, err = store.Get(id)
	require.Equal(t, sql.ErrNoRows, err)

}

// TestSetAddress
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, res.Address)
}

// TestSetStatus
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set status
	err = store.SetStatus(id, "delivered")
	require.NoError(t, err)

	// check
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Status, "delivered")
}

// TestGetByClient

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client
	// add
	for i := 0; i < len(parcels); i++ {

		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	// get by client

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {

		id := parcel.Number
		assert.NotEqual(t, 0, parcelMap[id])
		assert.Equal(t, parcelMap[id], parcel)

	}
}
