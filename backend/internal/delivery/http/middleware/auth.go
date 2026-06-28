package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"wishlist-go/internal/usecase/account"

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

func validateTelegramAuthData(rawAuthData string, hash string, botToken string) bool {
	rawAuthData = strings.TrimSpace(rawAuthData)
	rawAuthData = strings.Trim(rawAuthData, "'\"")
	hash = strings.TrimSpace(hash)
	hash = strings.Trim(hash, "'\"")

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

// tmaHeaderRe вытаскивает initData из заголовка "Authorization: tma <initData>".
var tmaHeaderRe = regexp.MustCompile(`^tma (.+)$`)

// parseAndValidateTelegram разбирает и проверяет подпись initData из заголовка Authorization.
// Возвращает данные и true только при валидной подписи и ненулевом user.id; ответ не пишет.
func parseAndValidateTelegram(authHeader, botToken string) (*TelegramAuthData, bool) {
	matches := tmaHeaderRe.FindStringSubmatch(authHeader)
	if matches == nil || len(matches) < 2 {
		return nil, false
	}
	rawAuthData := matches[1]
	values, err := url.ParseQuery(rawAuthData)
	if err != nil {
		return nil, false
	}

	var authData TelegramAuthData
	authData.QueryID = values.Get("query_id")
	authData.Hash = values.Get("hash")
	if authDateStr := values.Get("auth_date"); authDateStr != "" {
		authData.AuthDate, _ = strconv.ParseInt(authDateStr, 10, 64)
	}
	if userStr := values.Get("user"); userStr != "" {
		if err := json.Unmarshal([]byte(userStr), &authData.User); err != nil {
			return nil, false
		}
	}

	if !validateTelegramAuthData(rawAuthData, authData.Hash, botToken) {
		return nil, false
	}
	if authData.User.ID == 0 {
		return nil, false
	}
	return &authData, true
}

// TelegramAuthMiddleware требует валидный initData; иначе 401 и прерывание.
func TelegramAuthMiddleware(botToken string, accountUC *account.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authData, ok := parseAndValidateTelegram(c.GetHeader("Authorization"), botToken)
		if !ok {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// добавляем данные в контекст
		c.Set("telegram_auth", authData)

		// лениво создаём аккаунт синхронно, до передачи управления хендлеру
		if err := accountUC.EnsureExists(c.Request.Context(), authData.User.ID); err != nil {
			log.Printf("auth: не удалось создать/проверить аккаунт %d: %v", authData.User.ID, err)
		}

		c.Next()
	}
}

// OptionalTelegramAuth кладёт telegram_auth в контекст, ЕСЛИ заголовок валиден, но
// никогда не прерывает запрос. Для публичных эндпоинтов (гостевой share), которым нужно
// отличать владельца / гостя / анонима, но пускать запрос и без авторизации.
func OptionalTelegramAuth(botToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authData, ok := parseAndValidateTelegram(c.GetHeader("Authorization"), botToken); ok {
			c.Set("telegram_auth", authData)
		}
		c.Next()
	}
}
