package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"
    "strings"
	"word_press/auth"
	"word_press/models"
	"word_press/services"
)

type PrayerHandler struct {
	PrayerService   *services.PrayerService
	PresenceService *services.PresenceService

	GetUser func(*http.Request) (*models.User, error)
	Render  func(http.ResponseWriter, RenderOptions)
}

func NewPrayerHandler(
	ps *services.PrayerService,
	prs *services.PresenceService,
) *PrayerHandler {
	return &PrayerHandler{
		PrayerService:   ps,
		PresenceService: prs,
		GetUser:         auth.GetCurrentUser,
		Render:          RenderWithOpts,
	}
}

func (h *PrayerHandler) Home(w http.ResponseWriter, r *http.Request) {

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Presence heartbeat
	h.PresenceService.Touch(user)

	// Get stats
	stats := h.PresenceService.GetStats()

	// Get total prayer time
	totalSeconds, err := h.PrayerService.GetTotalTime(user.ID)
	if err != nil {
		log.Println("prayer total error:", err)
		totalSeconds = 0
	}

	greeting := getGreeting(user.Username)

	h.Render(w, RenderOptions{
		Page:   "prayer-home.html",
		Header: "home-header.html",
		Footer: "default",
		Data: map[string]interface{}{
			"Greeting":    greeting,
			"User":        user,
			"TotalTime":   formatDuration(totalSeconds),
			"ActiveUsers": stats.ActiveUsers,
			"Countries":   stats.Countries,
			"Cities":      stats.Cities,
		},
	})
}

func (h *PrayerHandler) Start(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.GetUser(r)
	log.Println(user)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID, err := h.PrayerService.Start(user.ID)
	log.Println(err)
	if err != nil {
		log.Println("start prayer error:", err)
		http.Error(w, "Failed to start prayer", http.StatusInternalServerError)
		return
	}

	// Ensure presence updated
	h.PresenceService.Touch(user)

	http.Redirect(w, r, "/prayer?started=1&session="+strconv.Itoa(sessionID), http.StatusSeeOther)
}

func (h *PrayerHandler) End(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	active, err := h.PrayerService.GetActive(user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch session", http.StatusInternalServerError)
		return
	}

	if active == nil {
		http.Error(w, "No active session", http.StatusBadRequest)
		return
	}

	err = h.PrayerService.End(active.ID)
	log.Println(err)
	if err != nil {
		http.Error(w, "Failed to end prayer", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/prayer?ended=1", http.StatusSeeOther)
}

func (h *PrayerHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.PresenceService.Touch(user)

	w.WriteHeader(http.StatusOK)
}

func (h *PrayerHandler) Fellow(w http.ResponseWriter, r *http.Request) {

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	h.PresenceService.Touch(user)

	users := h.PresenceService.GetUsers()

	h.Render(w, RenderOptions{
		Page:   "fellow-prayers.html",
		Header: "home-header.html",
		Footer: "default",
		Data: map[string]interface{}{
			"Users": users,
		},
	})
}

func (h *PrayerHandler) VisionNotes(w http.ResponseWriter, r *http.Request) {

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Presence tracking
	h.PresenceService.Touch(user)

	notes, err := models.GetVisionNotesByUser(user.ID)
	if err != nil {
		log.Println("vision notes error:", err)
		http.Error(w, "Failed to load notes", http.StatusInternalServerError)
		return
	}

	h.Render(w, RenderOptions{
		Page:   "vision-notes.html",
		Header: "home-header.html",
		Footer: "default",
		Data: map[string]interface{}{
			"Notes": notes,
		},
	})
}

func (h *PrayerHandler) CreateVisionNote(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))

	if content == "" {
		http.Redirect(w, r, "/vision-notes?error=empty", http.StatusSeeOther)
		return
	}

	err = models.CreateVisionNote(user.ID, content)
	if err != nil {
		log.Println("create vision note error:", err)
		http.Redirect(w, r, "/vision-notes?error=1", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/vision-notes?success=1", http.StatusSeeOther)
}

func (h *PrayerHandler) PrayerList(w http.ResponseWriter, r *http.Request) {

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Presence heartbeat
	h.PresenceService.Touch(user)

	prayers, err := models.GetPrayerItemsByUser(user.ID)
	if err != nil {
		log.Println("prayer list error:", err)
		http.Error(w, "Failed to load prayer list", http.StatusInternalServerError)
		return
	}

	h.Render(w, RenderOptions{
		Page:   "prayer-list.html",
		Header: "home-header.html",
		Footer: "default",
		Data: map[string]interface{}{
			"Prayers": prayers,
		},
	})
}

func (h *PrayerHandler) CreatePrayer(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	// 🔒 Validation
	if title == "" || content == "" {
		http.Redirect(w, r, "/prayer/list?error=empty", http.StatusSeeOther)
		return
	}

	if len(title) > 200 || len(content) > 5000 {
		http.Redirect(w, r, "/prayer/list?error=toolong", http.StatusSeeOther)
		return
	}

	err = models.CreatePrayerItem(user.ID, title, content)
	if err != nil {
		log.Println("create prayer error:", err)
		http.Redirect(w, r, "/prayer/list?error=1", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/prayer/list?success=1", http.StatusSeeOther)
}

func (h *PrayerHandler) DeletePrayer(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.GetUser(r)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse ID
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Fetch prayer item
	prayer, err := models.GetPrayerItemByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch prayer", http.StatusInternalServerError)
		return
	}
	if prayer == nil {
		http.Error(w, "Prayer not found", http.StatusNotFound)
		return
	}

	// 🔒 Ownership check (CRITICAL)
	if prayer.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete
	if err := models.DeletePrayerItem(id); err != nil {
		http.Error(w, "Failed to delete prayer", http.StatusInternalServerError)
		return
	}

	// Redirect back
	http.Redirect(w, r, "/prayer/list?deleted=1", http.StatusSeeOther)
}

func getGreeting(name string) string {
	hour := time.Now().Hour()

	switch {
	case hour < 12:
		return "Good Morning " + name
	case hour < 18:
		return "Good Afternoon " + name
	default:
		return "Good Evening " + name
	}
}

func formatDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60

	return strconv.Itoa(h) + "h " + strconv.Itoa(m) + "m"
}