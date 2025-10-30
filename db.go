package main

import (
	"context"
	"errors"
	"fmt"
	"nyooom/logging"
	"os"
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go"
)

type BasicDB interface {
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	SetHash(ctx context.Context, key string, values map[string]string) error
	GetHash(ctx context.Context, key string) (map[string]string, error)
	DeleteHash(ctx context.Context, key string) error
}

type ValkeyDB struct {
	db     valkey.Client
	prefix string
}

type AdvancedDB interface {
	GetVersion(ctx context.Context) (string, error)
	SetVersion(ctx context.Context, version string) error
	SetLink(ctx context.Context, link Link) error
	GetLink(ctx context.Context, linkID string) (Link, error)
	DeleteLink(ctx context.Context, linkID string) error
}

type DB struct {
	basicDB BasicDB
}

// Basic DB functions to have more complex DBs implement

func SetupDB() AdvancedDB {
	// Read arguments
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic("Failed to parse DB_PORT: " + err.Error())
	}
	dbAddress := os.Getenv("DB_ADDRESS")
	if dbAddress == "" {
		panic("DB_ADDRESS is not set")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		logging.Println("DB_PASSWORD is not set")
	}

	// Connect to Valkey
	dbURL := fmt.Sprintf("%s:%d", dbAddress, dbPort)
	dbConnectionOptions := valkey.ClientOption{
		InitAddress: []string{dbURL},
		Password:    dbPassword,
	}
	dbClient, err := valkey.NewClient(dbConnectionOptions)
	if err != nil {
		panic("Failed to connect to Valkey: " + err.Error())
	}

	// Save DB
	db := DB{
		basicDB: &ValkeyDB{
			db:     dbClient,
			prefix: "Nyooom:",
		},
	}

	// Check if DB is up to date
	db.versioning()

	return db
}

func (db DB) versioning() {
	expectedDBVersion := "1"
	currentVersion, err := db.GetVersion(context.Background())
	if err != nil || currentVersion != expectedDBVersion {
		logging.PrintErrStr("Failed to get current version: " + err.Error() + ". Attempting to set version to " + expectedDBVersion)
		db.SetVersion(context.Background(), expectedDBVersion)
		return
	}
}

func (db *ValkeyDB) Get(ctx context.Context, key string) (string, error) {
	value, err := db.db.Do(ctx, db.db.B().Get().Key(db.prefix+key).Build()).ToString()
	if err != nil {
		return "", errors.New("Could not get key " + key + ": " + err.Error())
	}
	return value, nil
}

// Sets a key in the database with some duration.
// If duration is 0, the key will be set with no expiration.
func (db *ValkeyDB) Set(ctx context.Context, key string, value string, duration time.Duration) error {
	if duration == 0 {
		err := db.db.Do(ctx, db.db.B().Set().Key(db.prefix+key).Value(value).Build()).Error()
		if err != nil {
			return errors.New("Could not set key " + key + " with no expiration: " + err.Error())
		}
	} else {
		err := db.db.Do(ctx, db.db.B().Set().Key(db.prefix+key).Value(value).Ex(duration).Build()).Error()
		if err != nil {
			return errors.New("Could not set key " + key + ": " + err.Error())
		}
	}
	return nil
}

func (db *ValkeyDB) Delete(ctx context.Context, key string) error {
	err := db.db.Do(ctx, db.db.B().Del().Key(db.prefix+key).Build()).Error()
	if err != nil {
		return errors.New("Could not delete key " + key + ": " + err.Error())
	}
	return nil
}

func (db *ValkeyDB) GetHash(ctx context.Context, key string) (map[string]string, error) {
	hash, err := db.db.Do(ctx, db.db.B().Hgetall().Key(key).Build()).AsStrMap()
	if err != nil {
		return nil, errors.New("Could not get hash for key " + key + ": " + err.Error())
	}
	return hash, nil
}

func (db *ValkeyDB) SetHash(ctx context.Context, key string, values map[string]string) error {
	hash := db.db.B().Hset().Key(key).FieldValue()
	for field, value := range values {
		hash = hash.FieldValue(field, value)
	}
	err := db.db.Do(ctx, hash.Build()).Error()
	if err != nil {
		return errors.New("Could not set hash for key " + key + ": " + err.Error())
	}
	return nil
}

func (db *ValkeyDB) DeleteHash(ctx context.Context, key string) error {
	err := db.db.Do(ctx, db.db.B().Del().Key(key).Build()).Error()
	if err != nil {
		return errors.New("Could not delete hash for key " + key + ": " + err.Error())
	}
	return nil
}

// Complex DB functions to have more complex DBs implement

func (db DB) GetVersion(ctx context.Context) (string, error) {
	version, err := db.basicDB.Get(ctx, "version")
	if err != nil {
		return "", errors.New("Could not get db version: " + err.Error())
	}
	return version, nil
}

func (db DB) SetVersion(ctx context.Context, version string) error {
	err := db.basicDB.Set(ctx, "version", version, 0)
	if err != nil {
		return errors.New("Could not set db version: " + err.Error())
	}
	return nil
}

func (db DB) GetLink(ctx context.Context, linkID string) (Link, error) {
	rawLink, err := db.basicDB.GetHash(ctx, linkID)
	if err != nil {
		return Link{}, errors.New("Could not get link " + linkID + ": " + err.Error())
	}
	clicks, err := strconv.Atoi(rawLink["clicks"])
	if err != nil {
		return Link{}, errors.New("Could not get clicks for link " + linkID + ": " + err.Error())
	}
	link := Link{
		ID:     linkID,
		Slug:   rawLink["slug"],
		URL:    rawLink["url"],
		Clicks: clicks,
	}
	return link, nil
}

func (db DB) SetLink(ctx context.Context, link Link) error {
	err := db.basicDB.SetHash(ctx, link.ID, map[string]string{
		"slug":   link.Slug,
		"url":    link.URL,
		"clicks": strconv.Itoa(link.Clicks),
	})
	if err != nil {
		return errors.New("Could not set link " + link.ID + ": " + err.Error())
	}
	return nil
}

func (db DB) DeleteLink(ctx context.Context, linkID string) error {
	err := db.basicDB.DeleteHash(ctx, linkID)
	if err != nil {
		return errors.New("Could not delete link " + linkID + ": " + err.Error())
	}
	return nil
}
