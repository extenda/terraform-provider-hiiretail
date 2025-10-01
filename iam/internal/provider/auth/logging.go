package auth

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"time"
)

// Logger provides structured logging for the auth package with automatic credential redaction
type Logger struct {
	*slog.Logger
	config LogConfig
}

// LogConfig configures logging behavior
type LogConfig struct {
	Level             slog.Level     // Log level (Debug, Info, Warn, Error)
	Format            LogFormat      // Output format (JSON, Text)
	Destination       LogDestination // Where to log (Console, File, Both)
	FilePath          string         // Log file path (if using File destination)
	RedactCredentials bool           // Whether to redact sensitive information
	AddSource         bool           // Whether to add source code location
}

// LogFormat represents the log output format
type LogFormat int

const (
	LogFormatJSON LogFormat = iota
	LogFormatText
)

// LogDestination represents where logs should be written
type LogDestination int

const (
	LogDestinationConsole LogDestination = iota
	LogDestinationFile
	LogDestinationBoth
)

// Sensitive data patterns for redaction
var (
	// OAuth2 tokens and secrets
	tokenPattern = regexp.MustCompile(`(?i)(access_token|refresh_token|client_secret|authorization|bearer)\s*[:=]\s*["\']?([a-zA-Z0-9\-._~+/]+=*)["\']?`)

	// Client credentials
	clientSecretPattern = regexp.MustCompile(`(?i)(client_secret|secret)\s*[:=]\s*["\']?([a-zA-Z0-9\-._~+/]+=*)["\']?`)

	// Authorization headers
	authHeaderPattern = regexp.MustCompile(`(?i)(authorization\s*:\s*bearer\s+)([a-zA-Z0-9\-._~+/]+=*)`)

	// URLs with embedded credentials
	urlCredentialPattern = regexp.MustCompile(`(https?://[^:/?#]+:)([^@/?#]+)(@[^/?#]+)`)
)

// Default logger instance
var defaultLogger *Logger

func init() {
	// Initialize default logger with safe defaults
	defaultLogger = NewLogger(LogConfig{
		Level:             slog.LevelInfo,
		Format:            LogFormatJSON,
		Destination:       LogDestinationConsole,
		RedactCredentials: true,
		AddSource:         false,
	})
}

// NewLogger creates a new structured logger with credential redaction
func NewLogger(config LogConfig) *Logger {
	var handler slog.Handler

	// Create handler options
	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	// Create base handler based on format
	switch config.Format {
	case LogFormatJSON:
		switch config.Destination {
		case LogDestinationConsole:
			handler = slog.NewJSONHandler(os.Stdout, opts)
		case LogDestinationFile:
			if config.FilePath != "" {
				file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
				if err != nil {
					// Fallback to console if file can't be opened
					handler = slog.NewJSONHandler(os.Stdout, opts)
				} else {
					handler = slog.NewJSONHandler(file, opts)
				}
			} else {
				handler = slog.NewJSONHandler(os.Stdout, opts)
			}
		case LogDestinationBoth:
			// For "both", we'll use console here and add file handling in middleware
			handler = slog.NewJSONHandler(os.Stdout, opts)
		}
	case LogFormatText:
		switch config.Destination {
		case LogDestinationConsole:
			handler = slog.NewTextHandler(os.Stdout, opts)
		case LogDestinationFile:
			if config.FilePath != "" {
				file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
				if err != nil {
					handler = slog.NewTextHandler(os.Stdout, opts)
				} else {
					handler = slog.NewTextHandler(file, opts)
				}
			} else {
				handler = slog.NewTextHandler(os.Stdout, opts)
			}
		case LogDestinationBoth:
			handler = slog.NewTextHandler(os.Stdout, opts)
		}
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	// Wrap handler with credential redaction if enabled
	if config.RedactCredentials {
		handler = &RedactingHandler{
			Handler: handler,
		}
	}

	return &Logger{
		Logger: slog.New(handler),
		config: config,
	}
}

// RedactingHandler wraps a slog.Handler to redact sensitive information
type RedactingHandler struct {
	slog.Handler
}

// RedactString redacts sensitive information from a string
func RedactString(s string) string {
	// Redact OAuth2 tokens and secrets
	s = tokenPattern.ReplaceAllStringFunc(s, func(match string) string {
		parts := tokenPattern.FindStringSubmatch(match)
		if len(parts) >= 3 {
			redacted := redactValue(parts[2])
			return parts[1] + "=" + redacted
		}
		return match
	})

	// Redact client secrets
	s = clientSecretPattern.ReplaceAllStringFunc(s, func(match string) string {
		parts := clientSecretPattern.FindStringSubmatch(match)
		if len(parts) >= 3 {
			redacted := redactValue(parts[2])
			return parts[1] + "=" + redacted
		}
		return match
	})

	// Redact authorization headers
	s = authHeaderPattern.ReplaceAllString(s, "${1}[REDACTED]")

	// Redact credentials in URLs
	s = urlCredentialPattern.ReplaceAllString(s, "${1}[REDACTED]${3}")

	return s
}

// redactValue redacts a credential value, showing first/last characters for debugging
func redactValue(value string) string {
	if len(value) <= 8 {
		return "[REDACTED]"
	}

	// Show first 4 and last 4 characters for debugging
	return value[:4] + "..." + value[len(value)-4:]
}

// Handle implements slog.Handler with credential redaction
func (h *RedactingHandler) Handle(ctx context.Context, record slog.Record) error {
	// Redact the message
	record.Message = RedactString(record.Message)

	// Redact attributes
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Value.Kind() == slog.KindString {
			attr.Value = slog.StringValue(RedactString(attr.Value.String()))
		}
		return true
	})

	return h.Handler.Handle(ctx, record)
}

// GetDefaultLogger returns the default logger instance
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// SetDefaultLogger sets the default logger instance
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

// WithContext adds contextual information to log entries
func (l *Logger) WithContext(ctx context.Context) *Logger {
	logger := l.Logger

	// Add request ID if present
	if requestID := ctx.Value("request_id"); requestID != nil {
		logger = logger.With("request_id", requestID)
	}

	// Add tenant ID if present
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		logger = logger.With("tenant_id", tenantID)
	}

	// Add client ID if present (but redact it)
	if clientID := ctx.Value("client_id"); clientID != nil {
		redactedClientID := redactValue(fmt.Sprint(clientID))
		logger = logger.With("client_id", redactedClientID)
	}

	return &Logger{
		Logger: logger,
		config: l.config,
	}
}

// WithFields adds structured fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	logger := l.Logger

	for key, value := range fields {
		// Redact sensitive field values using existing function from errors.go
		if l.config.RedactCredentials && isSensitiveField(key) {
			if strValue, ok := value.(string); ok {
				value = RedactString(strValue)
			}
		}
		logger = logger.With(key, value)
	}

	return &Logger{
		Logger: logger,
		config: l.config,
	}
}

// Convenience methods for different log levels with automatic credential redaction

// Debug logs debug information
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(RedactString(fmt.Sprintf(msg, args...)))
}

// Info logs informational messages
func Info(msg string, args ...interface{}) {
	defaultLogger.Info(RedactString(fmt.Sprintf(msg, args...)))
}

// Warn logs warning messages
func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(RedactString(fmt.Sprintf(msg, args...)))
}

// Error logs error messages
func Error(msg string, args ...interface{}) {
	defaultLogger.Error(RedactString(fmt.Sprintf(msg, args...)))
}

// DebugWithFields logs debug information with structured fields
func DebugWithFields(msg string, fields map[string]interface{}) {
	defaultLogger.WithFields(fields).Debug(RedactString(msg))
}

// InfoWithFields logs informational messages with structured fields
func InfoWithFields(msg string, fields map[string]interface{}) {
	defaultLogger.WithFields(fields).Info(RedactString(msg))
}

// WarnWithFields logs warning messages with structured fields
func WarnWithFields(msg string, fields map[string]interface{}) {
	defaultLogger.WithFields(fields).Warn(RedactString(msg))
}

// ErrorWithFields logs error messages with structured fields
func ErrorWithFields(msg string, fields map[string]interface{}) {
	defaultLogger.WithFields(fields).Error(RedactString(msg))
}

// OAuth2-specific logging functions

// LogTokenAcquisition logs OAuth2 token acquisition attempts
func LogTokenAcquisition(tenantID, clientID string, scopes []string, startTime time.Time) {
	duration := time.Since(startTime)
	InfoWithFields("OAuth2 token acquisition started", map[string]interface{}{
		"component":   "oauth2",
		"operation":   "token_acquisition",
		"tenant_id":   tenantID,
		"client_id":   clientID,
		"scopes":      scopes,
		"duration_ms": duration.Milliseconds(),
	})
}

// LogTokenSuccess logs successful OAuth2 token acquisition
func LogTokenSuccess(tenantID, clientID string, expiresIn int, startTime time.Time) {
	duration := time.Since(startTime)
	InfoWithFields("OAuth2 token acquired successfully", map[string]interface{}{
		"component":   "oauth2",
		"operation":   "token_acquisition",
		"status":      "success",
		"tenant_id":   tenantID,
		"client_id":   clientID,
		"expires_in":  expiresIn,
		"duration_ms": duration.Milliseconds(),
	})
}

// LogTokenError logs OAuth2 token acquisition errors
func LogTokenError(tenantID, clientID string, err error, startTime time.Time) {
	duration := time.Since(startTime)
	ErrorWithFields("OAuth2 token acquisition failed", map[string]interface{}{
		"component":   "oauth2",
		"operation":   "token_acquisition",
		"status":      "error",
		"tenant_id":   tenantID,
		"client_id":   clientID,
		"error":       err.Error(),
		"duration_ms": duration.Milliseconds(),
	})
}

// LogTokenRefresh logs OAuth2 token refresh attempts
func LogTokenRefresh(tenantID, clientID string, reason string, startTime time.Time) {
	InfoWithFields("OAuth2 token refresh initiated", map[string]interface{}{
		"component":  "oauth2",
		"operation":  "token_refresh",
		"tenant_id":  tenantID,
		"client_id":  clientID,
		"reason":     reason,
		"started_at": startTime.Format(time.RFC3339),
	})
}

// LogAPIRequest logs outgoing API requests
func LogAPIRequest(method, url string, tenantID string, startTime time.Time) {
	InfoWithFields("API request initiated", map[string]interface{}{
		"component":  "api_client",
		"operation":  "http_request",
		"method":     method,
		"url":        RedactString(url), // Redact any credentials in URL
		"tenant_id":  tenantID,
		"started_at": startTime.Format(time.RFC3339),
	})
}

// LogAPIResponse logs API response information
func LogAPIResponse(method, url string, statusCode int, tenantID string, startTime time.Time) {
	duration := time.Since(startTime)
	level := slog.LevelInfo
	status := "success"

	if statusCode >= 400 {
		level = slog.LevelWarn
		status = "client_error"
	}
	if statusCode >= 500 {
		level = slog.LevelError
		status = "server_error"
	}

	defaultLogger.WithFields(map[string]interface{}{
		"component":   "api_client",
		"operation":   "http_response",
		"method":      method,
		"url":         RedactString(url),
		"status_code": statusCode,
		"status":      status,
		"tenant_id":   tenantID,
		"duration_ms": duration.Milliseconds(),
	}).Log(context.Background(), level, "API request completed")
}

// LogDiscoveryAttempt logs OAuth2 endpoint discovery attempts
func LogDiscoveryAttempt(tenantID, environment string, startTime time.Time) {
	InfoWithFields("OAuth2 endpoint discovery started", map[string]interface{}{
		"component":   "discovery",
		"operation":   "endpoint_discovery",
		"tenant_id":   tenantID,
		"environment": environment,
		"started_at":  startTime.Format(time.RFC3339),
	})
}

// LogDiscoverySuccess logs successful endpoint discovery
func LogDiscoverySuccess(tenantID, authURL, apiURL string, startTime time.Time) {
	duration := time.Since(startTime)
	InfoWithFields("OAuth2 endpoint discovery completed", map[string]interface{}{
		"component":   "discovery",
		"operation":   "endpoint_discovery",
		"status":      "success",
		"tenant_id":   tenantID,
		"auth_url":    RedactString(authURL),
		"api_url":     RedactString(apiURL),
		"duration_ms": duration.Milliseconds(),
	})
}

// LogConfigValidation logs configuration validation attempts
func LogConfigValidation(tenantID string, validationResult bool, errors []string) {
	level := slog.LevelInfo
	status := "success"

	if !validationResult {
		level = slog.LevelError
		status = "failed"
	}

	fields := map[string]interface{}{
		"component": "config",
		"operation": "validation",
		"status":    status,
		"tenant_id": tenantID,
	}

	if len(errors) > 0 {
		fields["validation_errors"] = errors
	}

	message := "Configuration validation completed"
	if !validationResult {
		message = "Configuration validation failed"
	}

	defaultLogger.WithFields(fields).Log(context.Background(), level, message)
}

// LogCacheOperation logs token cache operations
func LogCacheOperation(operation, tenantID string, hit bool) {
	status := "miss"
	if hit {
		status = "hit"
	}

	InfoWithFields("Token cache operation", map[string]interface{}{
		"component": "cache",
		"operation": operation,
		"status":    status,
		"tenant_id": tenantID,
	})
}

// SetLogLevel updates the default logger's log level
func SetLogLevel(level slog.Level) {
	config := defaultLogger.config
	config.Level = level
	defaultLogger = NewLogger(config)
}

// AddLoggerToContext adds a logger instance to the context
func AddLoggerToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

// GetLoggerFromContext retrieves a logger from the context, or returns default
func GetLoggerFromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value("logger").(*Logger); ok {
		return logger
	}
	return defaultLogger
}

// LogPanicRecovery logs panic recovery information
func LogPanicRecovery(recovered interface{}, tenantID string) {
	// Get stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	ErrorWithFields("Panic recovered in auth package", map[string]interface{}{
		"component":   "auth",
		"operation":   "panic_recovery",
		"panic_value": recovered,
		"tenant_id":   tenantID,
		"stack_trace": stack,
	})
}
