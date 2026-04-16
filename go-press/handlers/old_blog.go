package handlers

import (
	"word_press/models"
	"word_press/auth"
	"word_press/utils"
	"net/http"
	"github.com/gorilla/mux"
    "strconv"
	"fmt"
	"time"
	"io"
	"os"
	"log"
	"image"
	"image/jpeg"
	"image/png"
)

func AdminPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts()

	user, _ := auth.GetCurrentUser(r) 

	if err != nil {
       log.Println("AdminPosts error:", err)
	   http.Error(w, err.Error(), http.StatusInternalServerError)
	   return
	}

	Render(w, "admin/posts.html", map[string]interface{}{
	    "SiteTitle": "Manage Posts",
	    "Posts": posts,
	    "User": user,
    })
}

func AdminCreatePost(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetCurrentUser(r)

	data := map[string]interface{}{
		"SiteTitle": "Create Post",
		"User": user,
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		slug := r.FormValue("slug")
		content := r.FormValue("content")
		status := r.FormValue("status")
		excerpt := r.FormValue("excerpt")

		if title == "" || content == "" {
			data["Error"] = "Title and Content required"
			data["Form"] = map[string]string{
				"title": title,
				"slug": slug,
				"content": content,
			}
			Render(w, "admin/create-post.html", data)
			return
		}

		if slug == "" {
			slug = utils.GenerateSlug(title)
		}

		var imagePath string

		file, handler, err := r.FormFile("featured_image")

		if err == nil {
			defer file.Close()

			os.MkdirAll("static/uploads", os.ModePerm)
			filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
			path := "static/uploads/" + filename

			dst, err := os.Create(path)
			if err == nil {
				
			}
			defer dst.Close()
			if handler.Size > 600*1024 {
                img, format, err := image.Decode(file)
				if err != nil {
					http.Error(w, "Invalid image", 400)
					return
				}
				switch format {
				case "jpeg", "jpg":
					jpeg.Encode(dst, img, &jpeg.Options{
						Quality: 70,
					})
				case "png":
					encoder := png.Encoder{
						CompressionLevel: png.BestCompression,
					}
					// do for both cases
					encoder.Encode(dst, img)
                default:
					http.Error(w, "Unsupported image type", 400)
					return
				}
			} else {
                io.Copy(dst, file)
			}
             
		    imagePath = "/uploads/" + filename 
		}

		err = models.CreatePost(
			title, 
			slug, 
			content, 
			status,
			user.ID,
			excerpt,
		    imagePath,
		)

		if err != nil {
			data["Error"] = err.Error()
			Render(w, "admin/create-post.html", data)
			return
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	Render(w, "admin/create-post.html", data)
}

func AdminEditPost(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetCurrentUser(r)

	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	

	post, err := models.GetPostByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"User": user,
		"Post": post,
	}

	if r.Method == http.MethodPost {

		title := r.FormValue("title")
		slug  := r.FormValue("slug")
		content := r.FormValue("content")
		status := r.FormValue("status")
		excerpt := r.FormValue("excerpt")

		// keep old image by default
		imagePath := post.FeaturedImage

		file, handler, err := r.FormFile("featured_image")

		if err == nil {
			defer file.Close()
			os.MkdirAll("static/uploads", os.ModePerm)

			filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
			path := "static/uploads/" + filename

			dst, err := os.Create(path)
			if err != nil {
				return
			}
			defer dst.Close()
			
			if handler.Size > 600*1024 {
				img, format, err := image.Decode(file)
				if err != nil {
					http.Error(w, "Invalid image", 400)
					return
				}
				switch format {
				case "jpeg", "jpg":
					jpeg.Encode(dst, img, &jpeg.Options{
						Quality: 70,
					})
				case "png":
					encoder := png.Encoder{
						CompressionLevel: png.BestCompression,
					}
					// do for both cases
					encoder.Encode(dst, img)
                default:
					http.Error(w, "Unsupported image type", 400)
					return
				}
			} else {
				// small file -> just copy
				io.Copy(dst, file)
			}

			imagePath = "/uploads/" + filename
		}

		log.Println("Editting image", err)
		
		err = models.UpdatePost(id,
			title,
			slug,
			content, 
			status, 
			excerpt, 
			user.ID, 
			imagePath)

		if err != nil {
			data["Error"] = err.Error()
			Render(w, "admin/edit-post.html", data)
			return
		}
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
		return
	}

	Render(w, "admin/edit-post.html", data)
}

func BlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	post, err := models.GetPostBySlug(slug)
	user, _ := auth.GetCurrentUser(r)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	comments, _ := models.GetCommentsByPost(post.ID)
	RenderWithOpts(w, RenderOptions{
		Page:   "blog-post.html",
		Header: "blog-header.html",
		Footer: "default",
		Data: map[string]interface{}{
			"Post":     post,
			"Comments": comments,
			"User":     user,
		},
	})
}

func AdminDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = models.DeletePost(id)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}
