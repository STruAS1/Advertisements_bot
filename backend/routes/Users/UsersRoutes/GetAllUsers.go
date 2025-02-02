package usersroutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"
	"tgbotBARAHOLKA/utilits"
	"time"

	"github.com/go-chi/chi/v5"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}

type MessageRequest struct {
	Message string `json:"message"`
}

type EditBalanceRequest struct {
	Balance int `json:"Balance"`
}

type BanRequest struct {
	Duration string `json:"duration,omitempty"`
	Reason   string `json:"reason"`
}

func parseCustomDuration(input string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+)([a-zA-Z]+)`) // Ищем число + суффикс (например, "2d")
	matches := re.FindAllStringSubmatch(input, -1)

	var totalDuration time.Duration
	for _, match := range matches {
		if len(match) != 3 {
			return 0, fmt.Errorf("invalid duration format")
		}
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration")
		}
		unit := match[2]

		switch unit {
		case "m":
			totalDuration += time.Duration(value) * time.Minute
		case "h":
			totalDuration += time.Duration(value) * time.Hour
		case "d":
			totalDuration += time.Duration(value) * 24 * time.Hour
		case "w":
			totalDuration += time.Duration(value) * 7 * 24 * time.Hour
		case "mo":
			totalDuration += time.Duration(value) * 30 * 24 * time.Hour
		case "y":
			totalDuration += time.Duration(value) * 365 * 24 * time.Hour
		default:
			return 0, fmt.Errorf("unsupported duration unit: %s", unit)
		}
	}

	return totalDuration, nil
}
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func GetAllUsers(r chi.Router) {
	r.Get("/Users", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		limitStr := queryParams.Get("limit")
		pageStr := queryParams.Get("page")

		var limit, page int
		var err error

		if limitStr == "" {
			limit = 20
		} else {
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit <= 0 {
				http.Error(w, "Invalid limit: must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		if pageStr == "" {
			page = 1
		} else {
			page, err = strconv.Atoi(pageStr)
			if err != nil || page <= 0 {
				http.Error(w, "Invalid page: must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		var users []models.User
		query := db.DB.Preload("Bans").Order("created_at desc").Limit(limit).Offset((page - 1) * limit)

		if err := query.Find(&users).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to fetch advertisements",
			})
			return
		}

		var totalRecords int64
		countQuery := db.DB.Model(&models.User{})

		if err := countQuery.Count(&totalRecords).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to count total records",
			})
			return
		}
		UsersWithCount := make([]map[string]interface{}, len(users))
		for i, user := range users {
			var CountOfAds int64
			db.DB.Model(&models.Advertisement{}).Where(&models.Advertisement{UserID: user.ID}).Count(&CountOfAds)
			isBanned, unbanDate := user.IsBanned()
			UsersWithCount[i] = map[string]interface{}{
				"ID":         user.ID,
				"Balance":    user.Balance,
				"CountOfAds": CountOfAds,
				"CreatedAt":  user.CreatedAt,
				"UserName":   user.Username,
				"FL":         user.FirstName + " " + user.LastName,
				"phone":      user.Phone,
				"city":       user.City,
				"TgID":       user.TelegramID,
				"isBanned":   isBanned,
				"unbanDate":  unbanDate,
			}
		}

		totalPages := (totalRecords + int64(limit) - 1) / int64(limit)

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data: map[string]interface{}{
				"advertisements": UsersWithCount,
				"currentPage":    page,
				"totalPages":     totalPages,
			},
		})
	})
	r.Put("/user", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		ID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		var Balance EditBalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&Balance); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		var User models.User
		db.DB.Where(models.User{ID: uint(ID)}).First(&User)
		if User.TelegramID == 0 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		if err := db.DB.Model(&models.User{}).
			Where("id = ?", uint(ID)).
			Updates(map[string]interface{}{
				"balance": Balance.Balance,
			}).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to update record",
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})

	})
	r.Delete("/user", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		ID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		db.DB.Delete(&models.User{}, ID)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
	r.Post("/Alert/All", func(w http.ResponseWriter, r *http.Request) {
		var Message MessageRequest
		if err := json.NewDecoder(r.Body).Decode(&Message); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		var Users []models.User
		db.DB.Find(&Users)
		for _, user := range Users {
			if user.TelegramID != 0 {
				go utilits.SendMessageToUser(Message.Message, int64(user.TelegramID))
			}
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
	r.Post("/Alert/user", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		ID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		var Message MessageRequest
		if err := json.NewDecoder(r.Body).Decode(&Message); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		var User models.User
		db.DB.Where(models.User{ID: uint(ID)}).First(&User)
		if User.TelegramID == 0 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		utilits.SendMessageToUser(Message.Message, int64(User.TelegramID))
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
}

func BanRoutes(r chi.Router) {
	r.Post("/ban", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")

		userID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}

		var req BanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "User not found"})
			return
		}

		msgText := fmt.Sprintf("🚫 Вас заблокировали!\nПричина: %s", req.Reason)

		if req.Duration != "" {
			duration, err := parseCustomDuration(req.Duration)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid duration format"})
				return
			}
			msgText += fmt.Sprintf("\nСрок: %s", req.Duration)
			models.BanUser(db.DB, user.ID, duration, req.Reason)
		} else {
			models.BanUserForever(db.DB, user.ID, req.Reason)
		}

		utilits.CheckAndKickUserFromChannel(user.TelegramID)
		utilits.SendMessageToUser(msgText, int64(user.TelegramID))
		writeJSON(w, http.StatusOK, map[string]string{"message": "User banned successfully"})
	})

	r.Post("/unban", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")

		userID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"message": "User not found"})
			return
		}

		if err := models.UnbanUser(db.DB, user.ID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to unban user"})
			return
		}
		utilits.SendMessageToUser("✅Вас разблокировали!", int64(user.TelegramID))

		writeJSON(w, http.StatusOK, map[string]string{"message": "User unbanned successfully"})
	})
}
