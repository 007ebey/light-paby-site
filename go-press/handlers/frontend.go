package handlers

import (
	"net/http"
    "word_press/models"
	"log"
	"strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
	success := r.URL.Query().Get("success")

	data := map[string]interface{}{
		"SiteTitle": "Pastor Aby & Pastor Smitha",
	}

	if success == "1" {
		data["Success"] = "Your message has been sent successfully!"
	}

	RenderWithOpts(w, RenderOptions{
		Page: "updated-index.html",
        Header: "home-header.html",
		Footer: "default",
		Data: data,
	})
}

func Home2(w http.ResponseWriter, r *http.Request) {
	Render(w, "index2.html", map[string]interface{}{
		"SiteTitle": "Cool Site: Layout 2",
	})
}

func SliderAdmin(w http.ResponseWriter, r *http.Request) {
	Render(w, "admin/create-slider.html", map[string]interface{}{
		"SiteTitle": "Admin: Slide Add",
	})
}

func AdminSliders(w http.ResponseWriter, r *http.Request) {
	Render(w, "admin/sliders.html", map[string]interface{}{
		"SiteTitle": "Admin: Slide",
	})
}

func Blog(w http.ResponseWriter, r *http.Request) {
    // We only want to show "published" posts to the public
    pageStr := r.URL.Query().Get("page")
	page := 1

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 10
	offset := (page - 1) * limit

	posts, err := models.GetPostsByStatusPaginated("published", limit, offset)
	if err != nil {
		log.Println("Blog fetch error:", err)
		http.Error(w, "Unable to load the blog at this time", http.StatusInternalServerError)
		return
	}
	allPosts, err := models.GetPostsByStatus("published")
	if err != nil {
		http.Error(w, "Count error", http.StatusInternalServerError)
		return 
	}
	total := len(allPosts)
	totalPages := (total + limit - 1) / limit

	var showAdminLinks bool
	user, ok := r.Context().Value("user").(*models.User)
	if ok && user != nil {
	    if user.Role == "admin" || user.Role == "editor" {
			showAdminLinks = true
		}
	}

	log.Println("Show admin link", showAdminLinks)

	log.Println("Number of blogs gotten", totalPages)

	data :=  map[string]interface{}{
		"SiteTitle":  "Blog",
		"Posts":      posts,
		"Page":       page,
		"TotalPages": totalPages,
		"ShowAdminLinks": showAdminLinks, 
	}
	RenderWithOpts(w, RenderOptions{
			Page:   "blog.html",
			Header: "home-header.html",
			Footer: "default",
			Data:   data,
	})
	return
}
