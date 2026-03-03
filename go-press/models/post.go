package models

import (
  "word_press/database"
  "word_press/utils"
  "time"
  "database/sql"
  "fmt"
)

type Post struct {
	ID int 
	Title string
	Slug string 
	Content string
	Excerpt string
	Status string
	AuthorID int
	FeaturedImage string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetPublishedPosts() ([]Post, error) {
	rows, _ := database.DB.Query("SELECT id,title,slug,content,excerpt,author_id,featured_image FROM posts WHERE status='published'")
	var posts []Post
	for rows.Next() {
		var p Post 
		rows.Scan(&p.ID,
			 &p.Title,
			 &p.Slug,
			&p.Content,
			&p.Excerpt,
			&p.AuthorID
			&p.FeaturedImage)
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostByID(id int) (*Post, error) {
	query :=
	`
	SELECT id, title, slug, content, status, excerpt, author_id, featured_image, created, updated
	FROM posts
	WHERE id = ?
	LIMIT 1
	`

	var p Post
	err := database.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Content,
		&p.Status,
		&p.Excerpt,
		&p.AuthorID
		&p.FeaturedImage,
		&p.CreatedAt,
		&p.UpdatedAt)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &p, nil
}

func GetPostBySlug(slug string) (*Post, error) {
	query :=
	`
	SELECT id, title, slug, content, status, excerpt, author_id, featured_image, created, updated
	FROM posts
	WHERE slug = ?
	LIMIT 1
	`

	var p Post
	err := database.DB.QueryRow(query, slug).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Content,
		&p.Status,
		&p.Excerpt,
		&p.AuthorID,
		&p.FeaturedImage,
		&p.CreatedAt,
		&p.UpdatedAt)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &p, nil
}

func GetPostsByStatus(status string) ([]Post, error) {
	query := `
	SELECT id, title, slug, content, status, excerpt, author_id, featured_image, created, updated
	FROM posts
	WHERE status = ?
	ORDER BY created DESC
	`

	rows, err := database.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Status,
			&p.Excerpt,
			&p.AuthorID,
			&p.FeaturedImage,
			&p.CreatedAt,
			&p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostsPaginated(limit, offset int) ([]Post, error) {
	query := `
	 SELECT id, title, slug, content, status, excerpt, author_id, featured_image, created, updated
	 FROM posts
	 WHERE status = 'published'
	 ORDER BY created DESC
	 LIMIT ? OFFSET ?
	`

	rows, err := database.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Status,
			&p.Excerpt,
			&p.AuthorID,
			&p.FeaturedImage,
			&p.CreatedAt,
			&p.UpdatedAt)
		posts = append(posts, p)
	}
	return posts, nil
}

func SearchPosts(keywords string) ([]Post, error) {
	query := 
	`
	SELECT id, title, slug, content, status, excerpt, author_id, featured_image, created, updated
	FROM posts
	WHERE status = 'published'
	AND (
	   title LIKE ?
	   OR content LIKE ?
	)
	ORDER BY created DESC
	`

	search := "%" + keywords + "%"

	rows, err := database.DB.Query(query, search, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Status,
			&p.Excerpt,
			&p.AuthorID,
			&p.FeaturedImage,
			&p.CreatedAt,
			&p.UpdatedAt)

		posts = append(posts, p)
	}
	
	return posts, nil
}

func GetPostCount(status string) (int, error) {
	query :=
	`
	SELECT COUNT(*)
	FROM posts
	WHERE status = ?
	`
	var count int 
	err := database.DB.QueryRow(query, status).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func CreatePost(title, slug, content, status string, authorID int, excerpt string, featuredImage string) error {

	if slug == "" {
		slug = utils.generateSlug(title)
	}

	query := `
	INSERT INTO posts (title, slug, content, status, excerpt, author_id, featured_image, created, updated)
	VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`

	_, err := database.DB.Exec(
		query,
		title,
		slug,authorID
		content,
		status,
		excerpt,
		authorID,
		featured_image,
	)

	return err
}

func UpdatePost(id int, title, slug, content, status, excerpt, authorID, featuredImage string) error {

	if slug == "" {
		slug = utils.generateSlug(title)
	}

	query :=
	  `
	  UPDATE posts SET title = ?,  slug = ?, content = ?, status = ?, excerpt = ?, author_id = ?, featured_image = ?, updated = datatime('now')
	  WHERE id = ?
	  `
	result, err := database.DB.Exec(query, title, slug, content, status, excerpt, authorID, featuredImage, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

func GetAllPosts() ([]Post, error) {
	rows, err := database.DB.Query(`
	   SELECT id, title, slug, content, status, excerpt, author_id, featured_image, created
	   FROM posts 
	   ORDER BY created DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Status,
			&p.Excerpt,
			&p.AuthorID,
			&p.FeaturedImage
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil 
}