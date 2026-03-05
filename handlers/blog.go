package handlers

import (
	"database/sql"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/loqun/fiber-server/helper"
	"github.com/loqun/fiber-server/sql_generated"
	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
)

// get the blog data from the request body
type BlogRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Author  string   `json:"author"`
	Tags    []string `json:"tags"`
}

func (h *Handler) Blog(c *fiber.Ctx) error {
	queries := sql_generated.New(h.db)

	// Pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 5 // Posts per page
	offset := (page - 1) * limit

	// Get paginated blogs
	blogs, err := queries.GetBlogsPaginated(c.Context(), sql_generated.GetBlogsPaginatedParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		log.Println("Error fetching blogs: database error")
		blogs = []sql_generated.GetBlogsPaginatedRow{}
	}

	// Get total count for pagination
	totalCount, err := queries.CountBlogs(c.Context())
	if err != nil {
		log.Println("Error counting blogs: database error")
		totalCount = 0
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Transform blogs for template
	type BlogData struct {
		ID        string
		Title     string
		Content   string
		Author    string
		CreatedAt time.Time
		ReadTime  int
		Tags      []string
	}

	blogData := make([]BlogData, len(blogs))
	for i, blog := range blogs {
		// Strip HTML tags from content for preview
		cleanContent := stripHTML(blog.Content)

		// Handle tag name
		tags := []string{}
		if blog.TagName.Valid {
			tags = append(tags, blog.TagName.String)
		}

		blogData[i] = BlogData{
			ID:        blog.ID,
			Title:     blog.Title,
			Content:   cleanContent,
			Author:    blog.Author,
			CreatedAt: blog.CreatedAt,
			ReadTime:  len(cleanContent)/200 + 1, // Rough estimate
			Tags:      tags,
		}
	}

	// check login state from session
	sess, _ := h.store.Get(c) // ignore error, on failure treat as not logged in
	isLogged := sess != nil && sess.Get("user") != nil

	data := fiber.Map{
		"Title":          "Personal Blog",
		"AuthorName":     "arif muftalib",
		"PageTitle":      "Personal Blog Entries",
		"Subtitle":       "Thoughts, experiences, and insights on technology and life.",
		"Blogs":          blogData,
		"CreatePostURL":  "/blog/form",
		"CreatePostText": "Write New Post",
		"IsLoggedIn":     isLogged,
		"CurrentPage":    page,
		"TotalPages":     totalPages,
		"HasPrevious":    page > 1,
		"HasNext":        page < totalPages,
		"PreviousPage":   page - 1,
		"NextPage":       page + 1,
	}

	return c.Render("blog_index", data)
}

func (h *Handler) GetAllBlog(c *fiber.Ctx) error {

	//fetch all the blog inside the database
	queries := sql_generated.New(h.db)
	blog, err := queries.GetAllBlog(c.Context())
	if err != nil {
		log.Println("Error fetching blogs:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch blogs"})
	}

	log.Println("Fetched blogs:", blog)

	//return as html page
	data := fiber.Map{"Blogs": blog, "Title": "All Blog Entries"}
	return c.Render("index.html", data)

}

func (h *Handler) StoreBlog(c *fiber.Ctx) error {

	// ensure user is authenticated; middleware should have run but double-check
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "session error"})
	}
	if sess.Get("user") == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var blogReq BlogRequest
	if err := c.BodyParser(&blogReq); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	//store the blog inside the database
	queries := sql_generated.New(h.db)
	createdAt := time.Now()

	// Generate ULID
	id := ulid.MustNew(ulid.Timestamp(createdAt), rand.New(rand.NewSource(createdAt.UnixNano())))

	var author string
	if lo.IsEmpty(blogReq.Author) {
		author = "Arif"
	} else {
		author = blogReq.Author
	}

	// Handle tag creation/assignment
	var tagID sql.NullString
	if len(blogReq.Tags) > 0 && blogReq.Tags[0] != "" {
		tagULID := ulid.MustNew(ulid.Timestamp(createdAt), rand.New(rand.NewSource(createdAt.UnixNano())))
		tag, err := queries.GetOrCreateTag(c.Context(), sql_generated.GetOrCreateTagParams{
			ID:   tagULID.String(),
			Name: blogReq.Tags[0], // Use first tag for now
		})
		if err != nil {
			log.Println("Error creating tag:", err)
			tagID = sql.NullString{Valid: false}
		} else {
			tagID = sql.NullString{String: tag.ID, Valid: true}
		}
	}

	blog, err := queries.StoreBlog(c.Context(), sql_generated.StoreBlogParams{
		ID:        id.String(),
		Title:     blogReq.Title,
		Content:   helper.SanitizeHTML(blogReq.Content),
		TagID:     tagID,
		Author:    author,
		CreatedAt: createdAt,
	})

	if err != nil {
		log.Println("Error storing blog:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store blog"})
	}

	return c.JSON(fiber.Map{"blog": blog})
}

func (h *Handler) BlogForm(c *fiber.Ctx) error {
	// middleware already ensures authenticated user; we can get the user if needed
	user := c.Locals("user")
	data := fiber.Map{
		"Title":             "Create Blog Entry",
		"AuthorName":        "arif muftalib",
		"PageTitle":         "Create Blog Entries",
		"TitlePlaceholder":  "Enter your blog title...",
		"ShowTags":          true,
		"TagsPlaceholder":   "Enter tags separated by commas...",
		"EditorPlaceholder": "Compose an epic...",
		"SubmitText":        "Submit",
		"SubmitURL":         "/store-blog",
		"SuccessMessage":    "Blog submitted successfully!",
		"User":              user,
	}

	return c.Render("blog_form", data)
}

// stripHTML removes HTML tags from content
func stripHTML(content string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	clean := re.ReplaceAllString(content, "")

	// Clean up extra whitespace
	clean = strings.TrimSpace(clean)
	clean = regexp.MustCompile(`\s+`).ReplaceAllString(clean, " ")

	return clean
}

func (h *Handler) GetBlog(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Looking for blog with ID")

	queries := sql_generated.New(h.db)
	blog, err := queries.GetBlogByID(c.Context(), id)
	if err != nil {
		log.Println("Error fetching blog by ID: not found")
		return c.Status(fiber.StatusNotFound).SendString("Blog not found")
	}

	data := fiber.Map{
		"Title":      blog.Title,
		"AuthorName": "arif muftalib",
		"Blog":       blog,
	}

	return c.Render("blog_detail", data)
}
