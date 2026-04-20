package database 


import ( 
	"database/sql" 
	_ "github.com/mattn/go-sqlite3" 
	"log"
 )


var DB *sql.DB

func InitDB() {
	db, err := sql.Open("sqlite3", "cms.db")
	if err != nil {
		log.Fatal("DB open error:", err)
	}
	DB = db

	// 🔥 enable FK enforcement
	if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	queries := []string{

		// USERS
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'subscriber',
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// POSTS
		`CREATE TABLE IF NOT EXISTS posts (
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
		);`,

		// COMMENTS
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			author TEXT,
			email TEXT,
			body TEXT,
			status TEXT,
			created DATETIME
		);`,

		// CONTACTS
		`CREATE TABLE IF NOT EXISTS contacts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			message TEXT NOT NULL,
			status TEXT DEFAULT 'new',
			ip_address TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// PRAYER SESSIONS
		`CREATE TABLE IF NOT EXISTS prayer_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			start_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			end_time DATETIME,
			duration INTEGER DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,

		`CREATE INDEX IF NOT EXISTS idx_prayer_sessions_user
		 ON prayer_sessions(user_id);`,

		`CREATE INDEX IF NOT EXISTS idx_prayer_sessions_active
		 ON prayer_sessions(user_id, end_time);`,

		// VISION NOTES
		`CREATE TABLE IF NOT EXISTS vision_notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,

		`CREATE INDEX IF NOT EXISTS idx_vision_notes_user
		 ON vision_notes(user_id);`,

		// PRAYER ITEMS
		`CREATE TABLE IF NOT EXISTS prayer_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,

		`CREATE INDEX IF NOT EXISTS idx_prayer_items_user_created
		 ON prayer_items(user_id, created_at);`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatal("DB schema error:", err)
		}
	}

	SeedAdminUser()
}