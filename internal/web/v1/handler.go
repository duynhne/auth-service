package v1

import (
	"errors"
	"net/http"

	"github.com/duynhne/auth-service/internal/core/domain"
	logicv1 "github.com/duynhne/auth-service/internal/logic/v1"
	"github.com/duynhne/auth-service/middleware"
	pkgzerolog "github.com/duynhne/pkg/logger/zerolog"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Handler groups HTTP handlers for the auth API v1.
// Dependencies are injected via the constructor — no global state.
type Handler struct {
	auth *logicv1.AuthService
}

// NewHandler creates a new Handler with the given AuthService.
func NewHandler(auth *logicv1.AuthService) *Handler {
	return &Handler{auth: auth}
}

// RegisterRoutes registers all auth API v1 routes on the given router group.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/login", h.Login)
	rg.POST("/auth/register", h.Register)
	rg.GET("/auth/me", h.GetMe)
}

// Login handles HTTP request for user login.
func (h *Handler) Login(c *gin.Context) {
	ctx, span := middleware.StartSpan(c.Request.Context(), "http.request", trace.WithAttributes(
		attribute.String("layer", "web"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.Request.URL.Path),
	))
	defer span.End()

	logger := pkgzerolog.FromContext(ctx)

	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.SetAttributes(attribute.Bool("request.valid", false))
		span.RecordError(err)
		logger.Error().Err(err).Msg("Invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.Bool("request.valid", true))

	// Call business logic layer
	response, err := h.auth.Login(ctx, req)
	if err != nil {
		span.RecordError(err)
		logger.Error().Err(err).Msg("Login failed")

		switch {
		case errors.Is(err, logicv1.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		case errors.Is(err, logicv1.ErrUserNotFound):
			// Don't reveal that user doesn't exist (security best practice)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		case errors.Is(err, logicv1.ErrPasswordExpired):
			c.JSON(http.StatusForbidden, gin.H{"error": "Password expired"})
		case errors.Is(err, logicv1.ErrAccountLocked):
			c.JSON(http.StatusForbidden, gin.H{"error": "Account locked"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	logger.Info().Str("user_id", response.User.ID).Msg("Login successful")
	c.JSON(http.StatusOK, response)
}

// Register handles HTTP request for user registration.
func (h *Handler) Register(c *gin.Context) {
	ctx, span := middleware.StartSpan(c.Request.Context(), "http.request", trace.WithAttributes(
		attribute.String("layer", "web"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.Request.URL.Path),
	))
	defer span.End()

	logger := pkgzerolog.FromContext(ctx)

	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.SetAttributes(attribute.Bool("request.valid", false))
		span.RecordError(err)
		logger.Error().Err(err).Msg("Invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.Bool("request.valid", true))

	// Call business logic layer
	response, err := h.auth.Register(ctx, req)
	if err != nil {
		span.RecordError(err)
		logger.Error().
			Err(err).
			Str("username", req.Username).
			Msg("Registration failed")

		switch {
		case errors.Is(err, logicv1.ErrUserExists):
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	logger.Info().Str("user_id", response.User.ID).Msg("Registration successful")
	c.JSON(http.StatusCreated, response)
}

// GetMe handles HTTP request to get current user from session token.
// GET /api/v1/auth/me
// Authorization: Bearer <token>
func (h *Handler) GetMe(c *gin.Context) {
	ctx, span := middleware.StartSpan(c.Request.Context(), "http.request", trace.WithAttributes(
		attribute.String("layer", "web"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.Request.URL.Path),
	))
	defer span.End()

	logger := pkgzerolog.FromContext(ctx)

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		span.SetAttributes(attribute.Bool("auth.present", false))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Expect "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		span.SetAttributes(attribute.Bool("auth.valid_format", false))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		return
	}
	token := authHeader[len(bearerPrefix):]

	span.SetAttributes(attribute.Bool("auth.present", true))

	// Lookup user by token
	user, err := h.auth.GetUserByToken(ctx, token)
	if err != nil {
		span.RecordError(err)
		logger.Warn().Err(err).Msg("Token lookup failed")

		switch {
		case errors.Is(err, logicv1.ErrSessionNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		case errors.Is(err, logicv1.ErrSessionExpired):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	logger.Info().Str("user_id", user.ID).Msg("Token validated")
	c.JSON(http.StatusOK, user)
}
