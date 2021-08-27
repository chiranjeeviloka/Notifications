package main

import (
	"github.com/FreedomCentral/central/cache"
	"github.com/FreedomCentral/central/queue"
	"gorm.io/gorm"
)

// Service struct holds all variables common to all handlers.
// That is why members have to be safe for concurrent use and do not cause race conditions!
type Service struct {
	// Many database coonection pools are safe for concurrent use. Check the docs of specific db package
	// you want to use.
	db    *gorm.DB
	users cache.Dict
	queue queue.Queue

	// These are sample members. Do not place anything that does not work with concurrency here.
	// If you MUST place anything that is not safe for concurrent use, then protect that
	// field against data races using go channels and/or types from `sync` package.
}

// TODO: Add other Server methods here.

// TODO: Do NOT add http handlers HERE! Instead add them to handler.go file.
