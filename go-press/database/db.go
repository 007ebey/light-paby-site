package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	db, _ := sql.Open("sqlite3", "cms.db")
	DB = db
	DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
	   id INTEGER PRIMARY KEY,
	   username TEXT UNIQUE NOT NULL,
	   email TEXT UNIQUE NOT NULL,
	   password TEXT UNIQUE NOT NULL,
	   role TEXT NOT NULL DEFAULT 'subscriber'
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
	`)
}