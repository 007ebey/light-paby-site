package handlers

import (
	"net/http"
	"time"
	"word_press/auth"
	"word_press/models"
	"io"
	"os"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {

	user, _ := auth.GetCurrentUser(r)

	data := map[string]interface{} {
		"SiteTitle": "Admin Dashboard",
		"User": user,
		"Now": time.Now(),
	}

	Render(w, "admin/dashboard.html", data)
}

func AdminCreateSlider(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"SiteTitle": "Create Slider",
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		subtitle := r.FormValue("subtitle")
		buttonText := r.FormValue("button_text")
		buttonURL := r.FormValue("button_url")
		status := r.FormValue("status")

		file, handler, err := r.FormFile("image")
		if err != nil {
			data["Error"] = "Image upload is required"
			Render(w, "admin/create-slider.html", data)
			return
		}

		defer file.Close()

		path := "static/uploads" + handler.Filename
		dst, err := os.Create(path)
		if err != nil {
			data["Error"] = "Failed to save image"
		    Render(w, "admin/create-slider.html", data)
			return
		}

		defer dst.Close()
			
		_, err = io.Copy(dst, file)

		if err != nil {
			data["Error"] = "Failed to write image"
			Render(w, "admin/create-slider.html", data)
			return
		}

		err = models.CreateSlider(
			title,
			subtitle,
			path,
			buttonText,
			buttonURL,
			status,
		)
		
		if err != nil {
			data["Error"] = "Failed to save to database"
			Render(w, "admin/create-slider.html", data)
			return
		}

		http.Redirect(w, r, "admin/sliders", http.StatusSeeOther)

		return
	}

	Render(w, "admin/create-slider.html", data)
}