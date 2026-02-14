package middleware

import (
	"VEDA95/open_board/api/internal/db/repository"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CORSMiddleware provides dynamic CORS handling that reads allowed origins
// from the database with caching to avoid DB queries on every request.
type CORSMiddleware struct {
	authSettingsRepo *repository.AuthSettingsRepository
	cacheTTL         time.Duration
	cachedOrigins    string
	cacheTime        time.Time
	mu               sync.RWMutex
}

// NewCORSMiddleware creates a new CORSMiddleware with the given repository
// and a default cache TTL of 30 seconds.
func NewCORSMiddleware(repo *repository.AuthSettingsRepository) *CORSMiddleware {
	return &CORSMiddleware{
		authSettingsRepo: repo,
		cacheTTL:         30 * time.Second,
		cachedOrigins:    "http://localhost:3000",
	}
}

// NewCORSMiddlewareWithTTL creates a new CORSMiddleware with a custom cache TTL.
func NewCORSMiddlewareWithTTL(repo *repository.AuthSettingsRepository, ttl time.Duration) *CORSMiddleware {
	return &CORSMiddleware{
		authSettingsRepo: repo,
		cacheTTL:         ttl,
		cachedOrigins:    "http://localhost:3000",
	}
}

// getOrigins retrieves the allowed origins, using cache if still valid.
func (m *CORSMiddleware) getOrigins() string {
	m.mu.RLock()
	if time.Since(m.cacheTime) < m.cacheTTL && m.cachedOrigins != "" {
		origins := m.cachedOrigins
		m.mu.RUnlock()
		return origins
	}
	m.mu.RUnlock()

	// Cache expired or empty, need to refresh
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock (another goroutine may have updated)
	if time.Since(m.cacheTime) < m.cacheTTL && m.cachedOrigins != "" {
		return m.cachedOrigins
	}

	// Fetch from database
	authSettings, err := m.authSettingsRepo.Find()
	if err == nil && len(authSettings.CORSDomain) > 0 {
		m.cachedOrigins = authSettings.CORSDomain
	} else {
		m.cachedOrigins = "http://localhost:3000"
	}
	m.cacheTime = time.Now()

	return m.cachedOrigins
}

// isOriginAllowed checks if the request origin is in the allowed origins list.
func (m *CORSMiddleware) isOriginAllowed(origin string, allowedOrigins string) bool {
	if allowedOrigins == "*" {
		return true
	}

	// Split allowed origins by comma and check each
	origins := strings.Split(allowedOrigins, ",")
	for _, allowed := range origins {
		allowed = strings.TrimSpace(allowed)
		if allowed == origin {
			return true
		}
	}

	return false
}

// Handler returns a Fiber handler that applies CORS headers dynamically.
func (m *CORSMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		// Skip CORS for same-origin requests
		if origin == "" {
			return c.Next()
		}

		allowedOrigins := m.getOrigins()

		// Handle preflight requests
		if c.Method() == fiber.MethodOptions {
			if m.isOriginAllowed(origin, allowedOrigins) {
				c.Set("Access-Control-Allow-Origin", origin)
				c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
				c.Set("Access-Control-Allow-Credentials", "true")
				c.Set("Access-Control-Max-Age", "86400")
			}
			return c.SendStatus(fiber.StatusNoContent)
		}

		// Handle actual requests
		if m.isOriginAllowed(origin, allowedOrigins) {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		}

		return c.Next()
	}
}
