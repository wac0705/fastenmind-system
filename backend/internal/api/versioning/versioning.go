package versioning

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// Version represents an API version
type Version struct {
	Major      int
	Minor      int
	Patch      int
	Deprecated bool
	Sunset     *time.Time
}

// APIVersionManager manages API versions
type APIVersionManager struct {
	versions       map[string]*Version
	routes         map[string]map[string]echo.HandlerFunc // version -> path -> handler
	defaultVersion string
	latestVersion  string
	config         *VersionConfig
}

// VersionConfig holds version configuration
type VersionConfig struct {
	// Version detection
	HeaderName       string // Default: "X-API-Version"
	QueryParamName   string // Default: "version"
	URLPathPrefix    bool   // Use /v1/resource format
	AcceptHeader     bool   // Use Accept: application/vnd.api+json;version=1
	
	// Deprecation settings
	DeprecationHeader    string // Default: "X-API-Deprecated"
	SunsetHeader         string // Default: "X-API-Sunset"
	MinSupportedVersion  string
	
	// Response headers
	VersionHeader        string // Default: "X-API-Version"
	SupportedHeader      string // Default: "X-API-Supported-Versions"
}

// NewAPIVersionManager creates a new API version manager
func NewAPIVersionManager(config *VersionConfig) *APIVersionManager {
	if config == nil {
		config = &VersionConfig{
			HeaderName:        "X-API-Version",
			QueryParamName:    "version",
			URLPathPrefix:     true,
			DeprecationHeader: "X-API-Deprecated",
			SunsetHeader:      "X-API-Sunset",
			VersionHeader:     "X-API-Version",
			SupportedHeader:   "X-API-Supported-Versions",
		}
	}
	
	return &APIVersionManager{
		versions: make(map[string]*Version),
		routes:   make(map[string]map[string]echo.HandlerFunc),
		config:   config,
	}
}

// RegisterVersion registers a new API version
func (vm *APIVersionManager) RegisterVersion(version string, v *Version) error {
	if _, exists := vm.versions[version]; exists {
		return fmt.Errorf("version %s already registered", version)
	}
	
	vm.versions[version] = v
	vm.routes[version] = make(map[string]echo.HandlerFunc)
	
	// Update latest version
	if vm.latestVersion == "" || vm.compareVersions(version, vm.latestVersion) > 0 {
		vm.latestVersion = version
	}
	
	// Set default if not set
	if vm.defaultVersion == "" {
		vm.defaultVersion = version
	}
	
	return nil
}

// VersionMiddleware returns middleware for API versioning
func (vm *APIVersionManager) VersionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Detect version
			version := vm.detectVersion(c)
			
			// Validate version
			v, exists := vm.versions[version]
			if !exists {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":             "Invalid API version",
					"requested_version": version,
					"supported_versions": vm.getSupportedVersions(),
				})
			}
			
			// Check if deprecated
			if v.Deprecated {
				c.Response().Header().Set(vm.config.DeprecationHeader, "true")
				if v.Sunset != nil {
					c.Response().Header().Set(vm.config.SunsetHeader, v.Sunset.Format(time.RFC3339))
				}
			}
			
			// Check minimum supported version
			if vm.config.MinSupportedVersion != "" && vm.compareVersions(version, vm.config.MinSupportedVersion) < 0 {
				return c.JSON(http.StatusGone, map[string]interface{}{
					"error":              "API version no longer supported",
					"requested_version":  version,
					"minimum_version":    vm.config.MinSupportedVersion,
					"supported_versions": vm.getSupportedVersions(),
				})
			}
			
			// Set version in context
			c.Set("api_version", version)
			c.Set("api_version_object", v)
			
			// Set response headers
			c.Response().Header().Set(vm.config.VersionHeader, version)
			c.Response().Header().Set(vm.config.SupportedHeader, strings.Join(vm.getSupportedVersions(), ", "))
			
			// Route to versioned handler
			path := c.Request().URL.Path
			if vm.config.URLPathPrefix {
				// Remove version prefix from path
				parts := strings.Split(path, "/")
				if len(parts) > 1 && strings.HasPrefix(parts[1], "v") {
					path = "/" + strings.Join(parts[2:], "/")
				}
			}
			
			// Check for versioned handler
			if handlers, ok := vm.routes[version]; ok {
				if handler, ok := handlers[path]; ok {
					return handler(c)
				}
			}
			
			// Fall through to default handler
			return next(c)
		}
	}
}

// detectVersion detects API version from request
func (vm *APIVersionManager) detectVersion(c echo.Context) string {
	// 1. Check URL path (e.g., /v1/resource)
	if vm.config.URLPathPrefix {
		parts := strings.Split(c.Request().URL.Path, "/")
		if len(parts) > 1 && strings.HasPrefix(parts[1], "v") {
			return parts[1][1:] // Remove 'v' prefix
		}
	}
	
	// 2. Check header
	if version := c.Request().Header.Get(vm.config.HeaderName); version != "" {
		return version
	}
	
	// 3. Check query parameter
	if version := c.QueryParam(vm.config.QueryParamName); version != "" {
		return version
	}
	
	// 4. Check Accept header
	if vm.config.AcceptHeader {
		accept := c.Request().Header.Get("Accept")
		if version := vm.parseAcceptHeader(accept); version != "" {
			return version
		}
	}
	
	// 5. Use default version
	return vm.defaultVersion
}

// parseAcceptHeader parses version from Accept header
func (vm *APIVersionManager) parseAcceptHeader(accept string) string {
	// Parse "application/vnd.api+json;version=1.0"
	parts := strings.Split(accept, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "version=") {
			return strings.TrimPrefix(part, "version=")
		}
	}
	return ""
}

// RegisterRoute registers a versioned route
func (vm *APIVersionManager) RegisterRoute(version, path string, handler echo.HandlerFunc) error {
	if _, exists := vm.versions[version]; !exists {
		return fmt.Errorf("version %s not registered", version)
	}
	
	vm.routes[version][path] = handler
	return nil
}

// Group creates a versioned route group
func (vm *APIVersionManager) Group(e *echo.Echo, version string) *echo.Group {
	if vm.config.URLPathPrefix {
		return e.Group("/v" + version)
	}
	return e.Group("")
}

// compareVersions compares two version strings
func (vm *APIVersionManager) compareVersions(v1, v2 string) int {
	// Parse semantic versions
	v1Parts := vm.parseVersion(v1)
	v2Parts := vm.parseVersion(v2)
	
	// Compare major
	if v1Parts[0] > v2Parts[0] {
		return 1
	} else if v1Parts[0] < v2Parts[0] {
		return -1
	}
	
	// Compare minor
	if v1Parts[1] > v2Parts[1] {
		return 1
	} else if v1Parts[1] < v2Parts[1] {
		return -1
	}
	
	// Compare patch
	if v1Parts[2] > v2Parts[2] {
		return 1
	} else if v1Parts[2] < v2Parts[2] {
		return -1
	}
	
	return 0
}

// parseVersion parses version string into components
func (vm *APIVersionManager) parseVersion(version string) [3]int {
	var parts [3]int
	
	// Handle simple numeric versions
	if v, err := strconv.Atoi(version); err == nil {
		parts[0] = v
		return parts
	}
	
	// Parse semantic version
	versionParts := strings.Split(version, ".")
	for i := 0; i < len(versionParts) && i < 3; i++ {
		if v, err := strconv.Atoi(versionParts[i]); err == nil {
			parts[i] = v
		}
	}
	
	return parts
}

// getSupportedVersions returns list of supported versions
func (vm *APIVersionManager) getSupportedVersions() []string {
	versions := make([]string, 0, len(vm.versions))
	for version, v := range vm.versions {
		if !v.Deprecated || v.Sunset == nil || v.Sunset.After(time.Now()) {
			versions = append(versions, version)
		}
	}
	return versions
}

// VersionedHandler wraps handlers with version-specific logic
type VersionedHandler struct {
	handlers map[string]echo.HandlerFunc
	fallback echo.HandlerFunc
}

// NewVersionedHandler creates a new versioned handler
func NewVersionedHandler() *VersionedHandler {
	return &VersionedHandler{
		handlers: make(map[string]echo.HandlerFunc),
	}
}

// AddVersion adds a version-specific handler
func (vh *VersionedHandler) AddVersion(version string, handler echo.HandlerFunc) {
	vh.handlers[version] = handler
}

// SetFallback sets the fallback handler
func (vh *VersionedHandler) SetFallback(handler echo.HandlerFunc) {
	vh.fallback = handler
}

// Handle returns the appropriate handler for the version
func (vh *VersionedHandler) Handle(c echo.Context) error {
	version := c.Get("api_version").(string)
	
	if handler, exists := vh.handlers[version]; exists {
		return handler(c)
	}
	
	if vh.fallback != nil {
		return vh.fallback(c)
	}
	
	return echo.NewHTTPError(http.StatusNotImplemented, "Version not implemented")
}

// VersionTransformer transforms responses between versions
type VersionTransformer interface {
	Transform(fromVersion, toVersion string, data interface{}) (interface{}, error)
}

// ResponseTransformer implements response transformation
type ResponseTransformer struct {
	transformers map[string]map[string]TransformFunc // fromVersion -> toVersion -> transform
}

type TransformFunc func(interface{}) (interface{}, error)

// NewResponseTransformer creates a new response transformer
func NewResponseTransformer() *ResponseTransformer {
	return &ResponseTransformer{
		transformers: make(map[string]map[string]TransformFunc),
	}
}

// RegisterTransform registers a transformation function
func (rt *ResponseTransformer) RegisterTransform(fromVersion, toVersion string, transform TransformFunc) {
	if _, exists := rt.transformers[fromVersion]; !exists {
		rt.transformers[fromVersion] = make(map[string]TransformFunc)
	}
	rt.transformers[fromVersion][toVersion] = transform
}

// Transform applies transformation between versions
func (rt *ResponseTransformer) Transform(fromVersion, toVersion string, data interface{}) (interface{}, error) {
	if fromVersion == toVersion {
		return data, nil
	}
	
	if transforms, exists := rt.transformers[fromVersion]; exists {
		if transform, exists := transforms[toVersion]; exists {
			return transform(data)
		}
	}
	
	return nil, fmt.Errorf("no transformation available from %s to %s", fromVersion, toVersion)
}

// VersionNegotiator handles content negotiation for API versions
type VersionNegotiator struct {
	supportedVersions []string
	defaultVersion    string
}

// NewVersionNegotiator creates a new version negotiator
func NewVersionNegotiator(supportedVersions []string, defaultVersion string) *VersionNegotiator {
	return &VersionNegotiator{
		supportedVersions: supportedVersions,
		defaultVersion:    defaultVersion,
	}
}

// Negotiate negotiates the best version based on client preferences
func (vn *VersionNegotiator) Negotiate(acceptHeader string) (string, error) {
	// Parse Accept header for version preferences
	// Implement content negotiation logic
	return vn.defaultVersion, nil
}

// VersionDeprecation handles version deprecation
type VersionDeprecation struct {
	version         string
	deprecatedAt    time.Time
	sunsetAt        *time.Time
	migrationGuide  string
	alternativeVersion string
}

// NewVersionDeprecation creates a new version deprecation
func NewVersionDeprecation(version string, sunsetAt *time.Time) *VersionDeprecation {
	return &VersionDeprecation{
		version:      version,
		deprecatedAt: time.Now(),
		sunsetAt:     sunsetAt,
	}
}

// SetMigrationGuide sets the migration guide URL
func (vd *VersionDeprecation) SetMigrationGuide(url string) {
	vd.migrationGuide = url
}

// SetAlternativeVersion sets the recommended alternative version
func (vd *VersionDeprecation) SetAlternativeVersion(version string) {
	vd.alternativeVersion = version
}

// GetDeprecationHeaders returns deprecation headers
func (vd *VersionDeprecation) GetDeprecationHeaders() map[string]string {
	headers := map[string]string{
		"X-API-Deprecated": "true",
		"X-API-Deprecated-At": vd.deprecatedAt.Format(time.RFC3339),
	}
	
	if vd.sunsetAt != nil {
		headers["X-API-Sunset"] = vd.sunsetAt.Format(time.RFC3339)
	}
	
	if vd.migrationGuide != "" {
		headers["X-API-Migration-Guide"] = vd.migrationGuide
	}
	
	if vd.alternativeVersion != "" {
		headers["X-API-Alternative-Version"] = vd.alternativeVersion
	}
	
	return headers
}

// Example transformations

// TransformV1ToV2 transforms response from v1 to v2 format
func TransformV1ToV2(data interface{}) (interface{}, error) {
	// Example: v1 uses "user_name", v2 uses "username"
	if m, ok := data.(map[string]interface{}); ok {
		if userName, exists := m["user_name"]; exists {
			m["username"] = userName
			delete(m, "user_name")
		}
	}
	return data, nil
}

// TransformV2ToV3 transforms response from v2 to v3 format
func TransformV2ToV3(data interface{}) (interface{}, error) {
	// Example: v3 adds new fields with defaults
	if m, ok := data.(map[string]interface{}); ok {
		if _, exists := m["api_version"]; !exists {
			m["api_version"] = "3.0.0"
		}
	}
	return data, nil
}