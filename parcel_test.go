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
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, parcel.Number, 0)

	// get
	storedParcel, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.Equal(t, parcel, storedParcel)

	// delete
	err = store.Delete(parcel.Number)
	require.NoError(t, err)
	parcel, err = store.Get(parcel.Number)
	require.Error(t, err)
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
	err = store.SetStatus(id, ParcelStatusDelivered)
	require.NoError(t, err)

	// check
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusDelivered, res.Status)
}

// TestGetByClient

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
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
		_, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		assert.EqualValues(t, parcelMap[parcel.Number], parcel)
	}
}
