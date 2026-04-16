package database

import (
	"log"
	"os"
	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
)

func SeedAdminUser() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (this is fine in prod)")
	}

	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASS")
	forceUpdate := os.Getenv("ADMIN_FORCE_UPDATE") == "true"

	if adminUser == "" || adminPass == "" {
		log.Println("Missing ADMIN_USER or ADMIN_PASS")
		return
	}

	var exists int
	err = DB.QueryRow(
		"SELECT COUNT(1) FROM users WHERE username = ?",
		adminUser,
	).Scan(&exists)

	if err != nil {
		log.Println("Seeder error:", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(adminPass),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Println("Password hash error:", err)
		return
	}

	if exists > 0 {
		if forceUpdate {
			_, err := DB.Exec(
				"UPDATE users SET password = ? WHERE username = ?",
				string(hashedPassword),
				adminUser,
			)
			if err != nil {
				log.Println("Password update error:", err)
				return
			}
			log.Println("Admin password updated")
		} else {
			log.Println("Admin user already exists (no update)")
		}
		return
	}

	// Create admin if not exists
	_, err = DB.Exec(`
		INSERT INTO users (username, email, password, role)
		VALUES (?, ?, ?, ?)     
	`,
		adminUser,
		"admin@example.com",
		string(hashedPassword),
		"admin",
	)

	if err != nil {
		log.Println("Seeder insert error:", err)
		return
	}

	log.Println("Default admin user created")
}