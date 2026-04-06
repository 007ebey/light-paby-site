package handlers

import (
	"log"
	"net/http"
	"strings"
	"word_press/models"
	"strconv"
)

func Contact(w http.ResponseWriter, r *http.Request) {

	log.Println("TEST HERE")

	// ---- GET → just render page ----
	if r.Method == http.MethodGet {

		success := r.URL.Query().Get("success")

		data := map[string]interface{}{
			"SiteTitle": "Pastor Aby & Pastor Smitha",
		}

		if success == "1" {
			data["Success"] = "Your message has been sent successfully!"
		}

		RenderWithOpts(w, RenderOptions{
			Page:   "updated-index.html",
			Header: "home-header.html",
			Footer: "default",
			Data:   data,
		})
		return
	}

	// ---- POST → process form ----
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form ONLY for POST
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Intent check (don’t trust random POSTs)
	if r.FormValue("form_type") != "contact" {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	message := strings.TrimSpace(r.FormValue("message"))

	// Validation
	if name == "" || email == "" || message == "" {
		RenderWithOpts(w, RenderOptions{
			Page:   "updated-index.html",
			Header: "home-header.html",
			Footer: "default",
			Data: map[string]interface{}{
				"Error":    "All fields are required",
				"FormData": map[string]string{"name": name, "email": email, "message": message},
				"SiteTitle": "Pastor Aby & Pastor Smitha",
			},
		})
		return
	}

	if !strings.Contains(email, "@") {
		RenderWithOpts(w, RenderOptions{
			Page:   "updated-index.html",
			Header: "home-header.html",
			Footer: "default",
			Data: map[string]interface{}{
				"Error":    "Invalid email address",
				"FormData": map[string]string{"name": name, "email": email, "message": message},
				"SiteTitle": "Pastor Aby & Pastor Smitha",
			},
		})
		return
	}

	ip := r.RemoteAddr

	err := models.CreateContact(name, email, message, ip)
	if err != nil {
		log.Println("Failed to save contact:", err)

		RenderWithOpts(w, RenderOptions{
			Page:   "updated-index.html",
			Header: "home-header.html",
			Footer: "default",
			Data: map[string]interface{}{
				"Error":    "Something went wrong. Please try again",
				"FormData": map[string]string{"name": name, "email": email, "message": message},
				"SiteTitle": "Pastor Aby & Pastor Smitha",
			},
		})
		return
	}

	log.Printf("Contact saved: %s %s\n", name, email)

	// ✅ CRITICAL FIX → Redirect instead of render
	http.Redirect(w, r, "/?success=1#contact-form", http.StatusSeeOther)
}

func ManageQueries(w http.ResponseWriter, r *http.Request) {

	// Get all contacts
	contacts, err := models.GetAllContacts()
	if err != nil {
		http.Error(w, "Failed to load contacts", http.StatusInternalServerError)
		return
	}

	// Prepare data
	data := map[string]interface{}{
		"SiteTitle": "Manage Queries",
		"Contacts":  contacts,
	}

	// Render template
	RenderWithOpts(w, RenderOptions{
		Page:   "admin/contact.html",
		Header: "contact-header.html",
		Footer: "default",
		Data:   data,
	})
}

func MarkContactRead(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "Missing contact ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = models.UpdateContactStatus(id, "read")
	if err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "Missing contact ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteContact(id)
	if err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
}