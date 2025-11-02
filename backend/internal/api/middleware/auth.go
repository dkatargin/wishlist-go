package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"wishlist-go/internal/config"
	"wishlist-go/internal/service"

	"github.com/gin-gonic/gin"
)

type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsPremium    bool   `json:"is_premium"`
}

type TelegramAuthData struct {
	QueryID  string       `json:"query_id"`
	User     TelegramUser `json:"user"`
	AuthDate int64        `json:"auth_date"`
	Hash     string       `json:"hash"`
}

func validateTelegramAuthData(rawAuthData string, hash string) bool {
	rawAuthData = strings.TrimSpace(rawAuthData)
	rawAuthData = strings.Trim(rawAuthData, "'\"")
	hash = strings.TrimSpace(hash)
	hash = strings.Trim(hash, "'\"")
	// Реализация проверки подписи данных Telegram
	botToken := config.Config.Telegram.BotToken

	values, err := url.ParseQuery(rawAuthData)
	if err != nil {
		return false
	}

	// Создаем data_check_string из отсортированных ключей (кроме hash)
	var keys []string
	for k := range values {
		if k != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var dataCheckParts []string
	for _, k := range keys {
		dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	dataCheckString := strings.Join(dataCheckParts, "\n")

	// Вычисляем HMAC
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	hmacKey := h.Sum(nil)

	finalHmac := hmac.New(sha256.New, hmacKey)
	finalHmac.Write([]byte(dataCheckString))
	finalHmacResult := hex.EncodeToString(finalHmac.Sum(nil))

	return finalHmacResult == hash
}

func TelegramAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		re := regexp.MustCompile(`^tma (.+)$`)
		matches := re.FindStringSubmatch(authHeader)
		if matches == nil || len(matches) < 2 {
			c.JSON(401, gin.H{"error": "wrong auth header"})
			c.Abort()
			return
		}

		rawAuthData := matches[1]
		values, err := url.ParseQuery(rawAuthData)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}
		// извлекаем данные
		var authData TelegramAuthData
		authData.QueryID = values.Get("query_id")
		authData.Hash = values.Get("hash")
		authData.User = TelegramUser{}

		// парсим auth_date
		if authDateStr := values.Get("auth_date"); authDateStr != "" {
			authData.AuthDate, _ = strconv.ParseInt(authDateStr, 10, 64)
		}

		// парсим user (это JSON в URL-encoded формате)
		if userStr := values.Get("user"); userStr != "" {
			if err := json.Unmarshal([]byte(userStr), &authData.User); err != nil {
				c.JSON(401, gin.H{"error": "invalid user data"})
				c.Abort()
				return
			}
		}
		if !validateTelegramAuthData(rawAuthData, authData.Hash) {
			c.JSON(401, gin.H{"error": "invalid Telegram auth data"})
			c.Abort()
			return
		}

		// добавляем данные в контекст
		c.Set("telegram_auth", &authData)

		// асинхронно создаем аккаунт, если его нет
		go func(authData TelegramAuthData) {
			account := service.NewAccountService()
			_, err := account.Get(authData.User.ID)
			if err != nil {
				if authData.User.ID == 0 {
					c.JSON(401, gin.H{"error": "invalid Telegram user ID"})
					c.Abort()
				} else {
					_, err := account.Create(authData.User.ID)
					if err != nil {
						fmt.Println("Failed to create account:", err)
					}
				}

			}
		}(authData)
		c.Next()
	}
}
