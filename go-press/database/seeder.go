package database

import (
	"log"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdminUser() {
	var exists int 

	err := DB.QueryRow(
		"SELECT COUNT(1) FROM users WHERE username = ?",
		"admin",
	).Scan(&exists)

	if err != nil {
		log.Println("Seeder error:", err)
		return
	}

	if exists > 0 {
		log.Println("Admin user already exists")
		return
	}

	password  := "admin123"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Println("Password hash error:", err)
		return
	}

	// insert user
	_, err = DB.Exec(`
	  INSERT INTO users (username, email, password, role)
	  VALUES (?, ?, ?, ?)     
	`,
      "admin",
	  "admin@example.com",
	  string(hashedPassword),
	  "admin",  
    )

	if err != nil {
		log.Println("Seeder insert error:", err)
		return
	}

	log.Println("Default admin user created")
	log.Println("Username: admin")
}