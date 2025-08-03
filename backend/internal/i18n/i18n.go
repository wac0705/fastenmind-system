package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*
var localesFS embed.FS

// I18n provides internationalization support
type I18n struct {
	translations map[string]map[string]interface{} // locale -> key -> value
	fallback     string
	supported    []language.Tag
	matcher      language.Matcher
	mu           sync.RWMutex
	config       *Config
}

// Config holds i18n configuration
type Config struct {
	// Default locale
	DefaultLocale string
	
	// Fallback locale
	FallbackLocale string
	
	// Supported locales
	SupportedLocales []string
	
	// File format
	Format string // "json", "yaml", "toml"
	
	// Translation directory
	LocalesDir string
	
	// Context key for locale
	ContextKey string
	
	// Accept-Language parsing
	ParseAcceptLanguage bool
	
	// Cookie settings
	CookieName string
	
	// Query parameter
	QueryParam string
}

// NewI18n creates a new i18n instance
func NewI18n(config *Config) (*I18n, error) {
	if config == nil {
		config = &Config{
			DefaultLocale:       "en",
			FallbackLocale:      "en",
			SupportedLocales:    []string{"en", "zh", "ja", "ko", "es", "fr", "de"},
			Format:              "json",
			LocalesDir:          "locales",
			ContextKey:          "locale",
			ParseAcceptLanguage: true,
			CookieName:          "locale",
			QueryParam:          "lang",
		}
	}
	
	i18n := &I18n{
		translations: make(map[string]map[string]interface{}),
		fallback:     config.FallbackLocale,
		config:       config,
	}
	
	// Load translations
	if err := i18n.loadTranslations(); err != nil {
		return nil, err
	}
	
	// Build language matcher
	i18n.buildMatcher()
	
	return i18n, nil
}

// loadTranslations loads all translation files
func (i *I18n) loadTranslations() error {
	// Load from embedded filesystem
	entries, err := localesFS.ReadDir("locales")
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		filename := entry.Name()
		locale := strings.TrimSuffix(filename, filepath.Ext(filename))
		
		// Check if locale is supported
		supported := false
		for _, sl := range i.config.SupportedLocales {
			if sl == locale {
				supported = true
				break
			}
		}
		
		if !supported {
			continue
		}
		
		// Load file
		data, err := localesFS.ReadFile(filepath.Join("locales", filename))
		if err != nil {
			return err
		}
		
		// Parse based on format
		translations := make(map[string]interface{})
		
		switch i.config.Format {
		case "json":
			if err := json.Unmarshal(data, &translations); err != nil {
				return err
			}
		case "yaml":
			if err := yaml.Unmarshal(data, &translations); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported format: %s", i.config.Format)
		}
		
		// Flatten nested translations
		flattened := make(map[string]interface{})
		i.flattenTranslations(translations, "", flattened)
		
		i.translations[locale] = flattened
	}
	
	return nil
}

// flattenTranslations flattens nested translation maps
func (i *I18n) flattenTranslations(data map[string]interface{}, prefix string, result map[string]interface{}) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		
		switch v := value.(type) {
		case map[string]interface{}:
			i.flattenTranslations(v, fullKey, result)
		default:
			result[fullKey] = value
		}
	}
}

// buildMatcher builds language matcher
func (i *I18n) buildMatcher() {
	tags := make([]language.Tag, 0, len(i.config.SupportedLocales))
	
	for _, locale := range i.config.SupportedLocales {
		tag, err := language.Parse(locale)
		if err == nil {
			tags = append(tags, tag)
		}
	}
	
	i.supported = tags
	i.matcher = language.NewMatcher(tags)
}

// Middleware returns i18n middleware for Echo
func (i *I18n) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			locale := i.detectLocale(c)
			
			// Set locale in context
			c.Set(i.config.ContextKey, locale)
			
			// Create localizer
			localizer := i.NewLocalizer(locale)
			c.Set("i18n", localizer)
			
			// Set Content-Language header
			c.Response().Header().Set("Content-Language", locale)
			
			return next(c)
		}
	}
}

// detectLocale detects locale from request
func (i *I18n) detectLocale(c echo.Context) string {
	// 1. Check query parameter
	if locale := c.QueryParam(i.config.QueryParam); locale != "" {
		if i.isSupported(locale) {
			return locale
		}
	}
	
	// 2. Check cookie
	if cookie, err := c.Cookie(i.config.CookieName); err == nil && cookie.Value != "" {
		if i.isSupported(cookie.Value) {
			return cookie.Value
		}
	}
	
	// 3. Check Accept-Language header
	if i.config.ParseAcceptLanguage {
		acceptLang := c.Request().Header.Get("Accept-Language")
		if acceptLang != "" {
			tags, _, err := language.ParseAcceptLanguage(acceptLang)
			if err == nil && len(tags) > 0 {
				match, _, _ := i.matcher.Match(tags...)
				return match.String()
			}
		}
	}
	
	// 4. Check user preference (if authenticated)
	if userLocale := c.Get("user_locale"); userLocale != nil {
		if locale, ok := userLocale.(string); ok && i.isSupported(locale) {
			return locale
		}
	}
	
	// 5. Use default
	return i.config.DefaultLocale
}

// isSupported checks if locale is supported
func (i *I18n) isSupported(locale string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	
	_, exists := i.translations[locale]
	return exists
}

// T translates a message
func (i *I18n) T(locale, key string, args ...interface{}) string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	
	// Get translations for locale
	translations, exists := i.translations[locale]
	if !exists {
		// Try fallback
		translations, exists = i.translations[i.fallback]
		if !exists {
			return key
		}
	}
	
	// Get translation
	value, exists := translations[key]
	if !exists {
		// Try fallback locale
		if locale != i.fallback {
			if fallbackTrans, ok := i.translations[i.fallback]; ok {
				if fallbackValue, ok := fallbackTrans[key]; ok {
					value = fallbackValue
				} else {
					return key
				}
			} else {
				return key
			}
		} else {
			return key
		}
	}
	
	// Format translation
	switch v := value.(type) {
	case string:
		if len(args) > 0 {
			return fmt.Sprintf(v, args...)
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Localizer provides localized translations
type Localizer struct {
	i18n   *I18n
	locale string
}

// NewLocalizer creates a new localizer
func (i *I18n) NewLocalizer(locale string) *Localizer {
	return &Localizer{
		i18n:   i,
		locale: locale,
	}
}

// T translates a message
func (l *Localizer) T(key string, args ...interface{}) string {
	return l.i18n.T(l.locale, key, args...)
}

// TP translates with pluralization
func (l *Localizer) TP(key string, count int, args ...interface{}) string {
	pluralKey := key
	
	// Simple pluralization rules
	switch l.locale {
	case "en":
		if count == 1 {
			pluralKey = key + ".one"
		} else {
			pluralKey = key + ".other"
		}
	case "zh", "ja", "ko":
		// No pluralization
		pluralKey = key + ".other"
	default:
		// Default pluralization
		if count == 1 {
			pluralKey = key + ".one"
		} else {
			pluralKey = key + ".other"
		}
	}
	
	// Try plural form first
	result := l.i18n.T(l.locale, pluralKey, append([]interface{}{count}, args...)...)
	if result != pluralKey {
		return result
	}
	
	// Fall back to base key
	return l.i18n.T(l.locale, key, append([]interface{}{count}, args...)...)
}

// TC translates with context
func (l *Localizer) TC(context, key string, args ...interface{}) string {
	contextKey := fmt.Sprintf("%s.%s", context, key)
	result := l.i18n.T(l.locale, contextKey, args...)
	if result != contextKey {
		return result
	}
	
	// Fall back to key without context
	return l.i18n.T(l.locale, key, args...)
}

// GetLocale returns the current locale
func (l *Localizer) GetLocale() string {
	return l.locale
}

// MessageBundle represents a collection of messages
type MessageBundle struct {
	Messages map[string]Message `json:"messages" yaml:"messages"`
}

// Message represents a translatable message
type Message struct {
	ID          string                 `json:"id" yaml:"id"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Default     string                 `json:"default" yaml:"default"`
	Translations map[string]Translation `json:"translations,omitempty" yaml:"translations,omitempty"`
}

// Translation represents a translated message
type Translation struct {
	Text   string            `json:"text" yaml:"text"`
	Plural map[string]string `json:"plural,omitempty" yaml:"plural,omitempty"`
}

// TemplateLocalizer provides template-based localization
type TemplateLocalizer struct {
	localizer *Localizer
	funcMap   template.FuncMap
}

// NewTemplateLocalizer creates a new template localizer
func NewTemplateLocalizer(localizer *Localizer) *TemplateLocalizer {
	tl := &TemplateLocalizer{
		localizer: localizer,
	}
	
	tl.funcMap = template.FuncMap{
		"t":  tl.localizer.T,
		"tp": tl.localizer.TP,
		"tc": tl.localizer.TC,
		"locale": func() string {
			return tl.localizer.GetLocale()
		},
		"formatDate": tl.formatDate,
		"formatTime": tl.formatTime,
		"formatNumber": tl.formatNumber,
		"formatCurrency": tl.formatCurrency,
	}
	
	return tl
}

// FuncMap returns template functions
func (tl *TemplateLocalizer) FuncMap() template.FuncMap {
	return tl.funcMap
}

// formatDate formats a date according to locale
func (tl *TemplateLocalizer) formatDate(date interface{}) string {
	// Implement locale-specific date formatting
	return fmt.Sprintf("%v", date)
}

// formatTime formats a time according to locale
func (tl *TemplateLocalizer) formatTime(time interface{}) string {
	// Implement locale-specific time formatting
	return fmt.Sprintf("%v", time)
}

// formatNumber formats a number according to locale
func (tl *TemplateLocalizer) formatNumber(number interface{}) string {
	// Implement locale-specific number formatting
	return fmt.Sprintf("%v", number)
}

// formatCurrency formats currency according to locale
func (tl *TemplateLocalizer) formatCurrency(amount interface{}, currency string) string {
	// Implement locale-specific currency formatting
	return fmt.Sprintf("%v %s", amount, currency)
}

// ValidationLocalizer provides localized validation messages
type ValidationLocalizer struct {
	localizer *Localizer
}

// NewValidationLocalizer creates a new validation localizer
func NewValidationLocalizer(localizer *Localizer) *ValidationLocalizer {
	return &ValidationLocalizer{
		localizer: localizer,
	}
}

// Required returns localized required field message
func (vl *ValidationLocalizer) Required(field string) string {
	return vl.localizer.T("validation.required", field)
}

// MinLength returns localized min length message
func (vl *ValidationLocalizer) MinLength(field string, min int) string {
	return vl.localizer.T("validation.min_length", field, min)
}

// MaxLength returns localized max length message
func (vl *ValidationLocalizer) MaxLength(field string, max int) string {
	return vl.localizer.T("validation.max_length", field, max)
}

// Email returns localized email validation message
func (vl *ValidationLocalizer) Email(field string) string {
	return vl.localizer.T("validation.email", field)
}

// Numeric returns localized numeric validation message
func (vl *ValidationLocalizer) Numeric(field string) string {
	return vl.localizer.T("validation.numeric", field)
}

// DateFormat returns localized date format validation message
func (vl *ValidationLocalizer) DateFormat(field, format string) string {
	return vl.localizer.T("validation.date_format", field, format)
}

// ErrorLocalizer provides localized error messages
type ErrorLocalizer struct {
	localizer *Localizer
}

// NewErrorLocalizer creates a new error localizer
func NewErrorLocalizer(localizer *Localizer) *ErrorLocalizer {
	return &ErrorLocalizer{
		localizer: localizer,
	}
}

// NotFound returns localized not found error
func (el *ErrorLocalizer) NotFound(resource string) string {
	return el.localizer.T("errors.not_found", resource)
}

// Unauthorized returns localized unauthorized error
func (el *ErrorLocalizer) Unauthorized() string {
	return el.localizer.T("errors.unauthorized")
}

// Forbidden returns localized forbidden error
func (el *ErrorLocalizer) Forbidden() string {
	return el.localizer.T("errors.forbidden")
}

// InternalError returns localized internal error
func (el *ErrorLocalizer) InternalError() string {
	return el.localizer.T("errors.internal_server_error")
}

// BadRequest returns localized bad request error
func (el *ErrorLocalizer) BadRequest(details string) string {
	return el.localizer.T("errors.bad_request", details)
}

// RateLimited returns localized rate limit error
func (el *ErrorLocalizer) RateLimited(retryAfter int) string {
	return el.localizer.T("errors.rate_limited", retryAfter)
}