package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB() {
	db, _ := sql.Open("sqlite3", "cms.db")
	DB = db
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
	   id INTEGER PRIMARY KEY AUTOINCREMENT,
	   username TEXT UNIQUE NOT NULL,
	   email TEXT UNIQUE NOT NULL,
	   password TEXT NOT NULL,
	   role TEXT NOT NULL DEFAULT 'subscriber',
	   created DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS posts (
	   id INTEGER PRIMARY KEY AUTOINCREMENT,
	   title TEXT NOT NULL,
	   slug TEXT UNIQUE NOT NULL,
	   content TEXT NOT NULL,
	   excerpt TEXT,
	   status TEXT NOT NULL DEFAULT 'draft',
	   author_id INTEGER,
	   featured_image TEXT,
	   created DATETIME DEFAULT CURRENT_TIMESTAMP,
	   updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS comments (
	   id INTEGER PRIMARY KEY AUTOINCREMENT,
	   post_id INTEGER,
	   author TEXT,
	   email TEXT,
	   body TEXT,
	   status TEXT,
	   created DATETIME
	);

	-- CONTACTS TABLE (new)
	CREATE TABLE IF NOT EXISTS contacts (
	   id INTEGER PRIMARY KEY AUTOINCREMENT,
	   name TEXT NOT NULL,
	   email TEXT NOT NULL,
	   message TEXT NOT NULL,
	   status TEXT DEFAULT 'new', -- new, read, replied, spam
	   ip_address TEXT,
	   created DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	`)

	log.Println(err)

	SeedAdminUser()
}