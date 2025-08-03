package servicemesh

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// IstioManager manages Istio service mesh configurations
type IstioManager struct {
	config          *IstioConfig
	virtualServices map[string]*VirtualService
	destinationRules map[string]*DestinationRule
	gateways        map[string]*Gateway
	policies        map[string]*SecurityPolicy
	mu              sync.RWMutex
}

// IstioConfig holds Istio configuration
type IstioConfig struct {
	Namespace          string
	IngressGateway     string
	TLSMode           string // "SIMPLE", "MUTUAL", "ISTIO_MUTUAL"
	mTLSEnabled       bool
	TracingEnabled    bool
	MetricsEnabled    bool
	AccessLogEnabled  bool
	CircuitBreaker    *CircuitBreakerConfig
	Retry             *RetryConfig
	Timeout           *TimeoutConfig
}

// VirtualService represents Istio VirtualService configuration
type VirtualService struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Namespace   string           `json:"namespace"`
	Hosts       []string         `json:"hosts"`
	Gateways    []string         `json:"gateways"`
	HTTPRoutes  []HTTPRoute      `json:"http_routes"`
	TCPRoutes   []TCPRoute       `json:"tcp_routes,omitempty"`
	TLSRoutes   []TLSRoute       `json:"tls_routes,omitempty"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// HTTPRoute represents HTTP routing rules
type HTTPRoute struct {
	Name        string            `json:"name"`
	Match       []HTTPMatchRequest `json:"match"`
	Route       []HTTPRouteDestination `json:"route"`
	Redirect    *HTTPRedirect     `json:"redirect,omitempty"`
	Rewrite     *HTTPRewrite      `json:"rewrite,omitempty"`
	Timeout     time.Duration     `json:"timeout,omitempty"`
	Retry       *HTTPRetry        `json:"retry,omitempty"`
	Fault       *HTTPFaultInjection `json:"fault,omitempty"`
	Mirror      *Destination      `json:"mirror,omitempty"`
	MirrorPercent *float32        `json:"mirror_percent,omitempty"`
	Headers     *Headers          `json:"headers,omitempty"`
}

// HTTPMatchRequest represents HTTP match conditions
type HTTPMatchRequest struct {
	Name        string            `json:"name,omitempty"`
	URI         *StringMatch      `json:"uri,omitempty"`
	Scheme      *StringMatch      `json:"scheme,omitempty"`
	Method      *StringMatch      `json:"method,omitempty"`
	Authority   *StringMatch      `json:"authority,omitempty"`
	Headers     map[string]*StringMatch `json:"headers,omitempty"`
	Port        uint32            `json:"port,omitempty"`
	SourceLabels map[string]string `json:"source_labels,omitempty"`
	Gateways    []string          `json:"gateways,omitempty"`
	QueryParams map[string]*StringMatch `json:"query_params,omitempty"`
}

// StringMatch represents string matching conditions
type StringMatch struct {
	Exact  string `json:"exact,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

// HTTPRouteDestination represents route destination
type HTTPRouteDestination struct {
	Destination *Destination `json:"destination"`
	Weight      int32        `json:"weight,omitempty"`
	Headers     *Headers     `json:"headers,omitempty"`
}

// Destination represents a destination service
type Destination struct {
	Host   string `json:"host"`
	Subset string `json:"subset,omitempty"`
	Port   *PortSelector `json:"port,omitempty"`
}

// PortSelector represents port selection
type PortSelector struct {
	Number uint32 `json:"number,omitempty"`
	Name   string `json:"name,omitempty"`
}

// HTTPRedirect represents HTTP redirect configuration
type HTTPRedirect struct {
	URI        string `json:"uri,omitempty"`
	Authority  string `json:"authority,omitempty"`
	Scheme     string `json:"scheme,omitempty"`
	RedirectCode uint32 `json:"redirect_code,omitempty"`
}

// HTTPRewrite represents HTTP rewrite configuration
type HTTPRewrite struct {
	URI       string `json:"uri,omitempty"`
	Authority string `json:"authority,omitempty"`
}

// HTTPRetry represents HTTP retry configuration
type HTTPRetry struct {
	Attempts      int32         `json:"attempts"`
	PerTryTimeout time.Duration `json:"per_try_timeout,omitempty"`
	RetryOn       string        `json:"retry_on,omitempty"`
	RetryRemoteLocalities bool  `json:"retry_remote_localities,omitempty"`
}

// HTTPFaultInjection represents fault injection configuration
type HTTPFaultInjection struct {
	Delay *HTTPFaultInjectionDelay `json:"delay,omitempty"`
	Abort *HTTPFaultInjectionAbort `json:"abort,omitempty"`
}

// HTTPFaultInjectionDelay represents delay fault injection
type HTTPFaultInjectionDelay struct {
	Percentage *Percent      `json:"percentage,omitempty"`
	FixedDelay time.Duration `json:"fixed_delay,omitempty"`
}

// HTTPFaultInjectionAbort represents abort fault injection
type HTTPFaultInjectionAbort struct {
	Percentage *Percent `json:"percentage,omitempty"`
	HTTPStatus int32    `json:"http_status,omitempty"`
}

// Percent represents percentage value
type Percent struct {
	Value float64 `json:"value"`
}

// Headers represents header manipulation
type Headers struct {
	Request  *HeaderOperations `json:"request,omitempty"`
	Response *HeaderOperations `json:"response,omitempty"`
}

// HeaderOperations represents header operations
type HeaderOperations struct {
	Set    map[string]string `json:"set,omitempty"`
	Add    map[string]string `json:"add,omitempty"`
	Remove []string          `json:"remove,omitempty"`
}

// TCPRoute represents TCP routing rules
type TCPRoute struct {
	Match []TCPRouteMatch `json:"match,omitempty"`
	Route []TCPRouteDestination `json:"route"`
}

// TCPRouteMatch represents TCP match conditions
type TCPRouteMatch struct {
	DestinationSubnets []string          `json:"destination_subnets,omitempty"`
	Port              uint32            `json:"port,omitempty"`
	SourceLabels      map[string]string `json:"source_labels,omitempty"`
	Gateways          []string          `json:"gateways,omitempty"`
}

// TCPRouteDestination represents TCP route destination
type TCPRouteDestination struct {
	Destination *Destination `json:"destination"`
	Weight      int32        `json:"weight,omitempty"`
}

// TLSRoute represents TLS routing rules
type TLSRoute struct {
	Match []TLSRouteMatch `json:"match"`
	Route []TCPRouteDestination `json:"route"`
}

// TLSRouteMatch represents TLS match conditions
type TLSRouteMatch struct {
	SNIHosts     []string          `json:"sni_hosts,omitempty"`
	DestinationSubnets []string    `json:"destination_subnets,omitempty"`
	Port         uint32            `json:"port,omitempty"`
	SourceLabels map[string]string `json:"source_labels,omitempty"`
	Gateways     []string          `json:"gateways,omitempty"`
}

// DestinationRule represents Istio DestinationRule configuration
type DestinationRule struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	Namespace       string           `json:"namespace"`
	Host            string           `json:"host"`
	TrafficPolicy   *TrafficPolicy   `json:"traffic_policy,omitempty"`
	Subsets         []Subset         `json:"subsets,omitempty"`
	Labels          map[string]string `json:"labels"`
	Annotations     map[string]string `json:"annotations"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// TrafficPolicy represents traffic policy configuration
type TrafficPolicy struct {
	LoadBalancer    *LoadBalancerSettings `json:"load_balancer,omitempty"`
	ConnectionPool  *ConnectionPoolSettings `json:"connection_pool,omitempty"`
	OutlierDetection *OutlierDetection    `json:"outlier_detection,omitempty"`
	TLS             *ClientTLSSettings   `json:"tls,omitempty"`
	PortLevelSettings []PortTrafficPolicy `json:"port_level_settings,omitempty"`
}

// LoadBalancerSettings represents load balancer configuration
type LoadBalancerSettings struct {
	Simple         string                     `json:"simple,omitempty"` // ROUND_ROBIN, LEAST_CONN, RANDOM, PASSTHROUGH
	ConsistentHash *ConsistentHashLB          `json:"consistent_hash,omitempty"`
	LocalityLbSetting *LocalityLoadBalancerSetting `json:"locality_lb_setting,omitempty"`
}

// ConsistentHashLB represents consistent hash load balancer
type ConsistentHashLB struct {
	HTTPCookieName string        `json:"http_cookie_name,omitempty"`
	HTTPCookieTTL  time.Duration `json:"http_cookie_ttl,omitempty"`
	HTTPHeaderName string        `json:"http_header_name,omitempty"`
	UseSourceIP    bool          `json:"use_source_ip,omitempty"`
	RingHash       *RingHashLB   `json:"ring_hash,omitempty"`
	MagLev         *MagLevLB     `json:"maglev,omitempty"`
}

// RingHashLB represents ring hash configuration
type RingHashLB struct {
	MinimumRingSize uint64 `json:"minimum_ring_size,omitempty"`
}

// MagLevLB represents MagLev configuration
type MagLevLB struct {
	TableSize uint64 `json:"table_size,omitempty"`
}

// LocalityLoadBalancerSetting represents locality-aware load balancing
type LocalityLoadBalancerSetting struct {
	Distribute []LocalityLoadBalancerSettingDistribute `json:"distribute,omitempty"`
	Failover   []LocalityLoadBalancerSettingFailover   `json:"failover,omitempty"`
	Enabled    *bool                                   `json:"enabled,omitempty"`
}

// LocalityLoadBalancerSettingDistribute represents distribution settings
type LocalityLoadBalancerSettingDistribute struct {
	From string            `json:"from,omitempty"`
	To   map[string]uint32 `json:"to,omitempty"`
}

// LocalityLoadBalancerSettingFailover represents failover settings
type LocalityLoadBalancerSettingFailover struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// ConnectionPoolSettings represents connection pool configuration
type ConnectionPoolSettings struct {
	TCP  *TCPSettings  `json:"tcp,omitempty"`
	HTTP *HTTPSettings `json:"http,omitempty"`
}

// TCPSettings represents TCP connection settings
type TCPSettings struct {
	MaxConnections int32         `json:"max_connections,omitempty"`
	ConnectTimeout time.Duration `json:"connect_timeout,omitempty"`
	TCPNoDelay     bool          `json:"tcp_no_delay,omitempty"`
	TCPKeepAlive   *TCPKeepAlive `json:"tcp_keep_alive,omitempty"`
}

// TCPKeepAlive represents TCP keep-alive settings
type TCPKeepAlive struct {
	Time     time.Duration `json:"time,omitempty"`
	Interval time.Duration `json:"interval,omitempty"`
	Probes   uint32        `json:"probes,omitempty"`
}

// HTTPSettings represents HTTP connection settings
type HTTPSettings struct {
	HTTP1MaxPendingRequests  int32         `json:"http1_max_pending_requests,omitempty"`
	HTTP2MaxRequests         int32         `json:"http2_max_requests,omitempty"`
	MaxRequestsPerConnection int32         `json:"max_requests_per_connection,omitempty"`
	MaxRetries               int32         `json:"max_retries,omitempty"`
	IdleTimeout              time.Duration `json:"idle_timeout,omitempty"`
	H2UpgradePolicy          string        `json:"h2_upgrade_policy,omitempty"` // UPGRADE, DO_NOT_UPGRADE
	UseClientProtocol        bool          `json:"use_client_protocol,omitempty"`
}

// OutlierDetection represents outlier detection configuration
type OutlierDetection struct {
	ConsecutiveGatewayErrors uint32        `json:"consecutive_gateway_errors,omitempty"`
	Consecutive5xxErrors     uint32        `json:"consecutive_5xx_errors,omitempty"`
	Interval                 time.Duration `json:"interval,omitempty"`
	BaseEjectionTime         time.Duration `json:"base_ejection_time,omitempty"`
	MaxEjectionPercent       int32         `json:"max_ejection_percent,omitempty"`
	MinHealthPercent         int32         `json:"min_health_percent,omitempty"`
	SplitExternalLocalOriginErrors bool    `json:"split_external_local_origin_errors,omitempty"`
}

// ClientTLSSettings represents client TLS configuration
type ClientTLSSettings struct {
	Mode              string   `json:"mode"` // DISABLE, SIMPLE, MUTUAL, ISTIO_MUTUAL
	ClientCertificate string   `json:"client_certificate,omitempty"`
	PrivateKey        string   `json:"private_key,omitempty"`
	CACertificates    string   `json:"ca_certificates,omitempty"`
	SubjectAltNames   []string `json:"subject_alt_names,omitempty"`
	SNI               string   `json:"sni,omitempty"`
	InsecureSkipVerify bool    `json:"insecure_skip_verify,omitempty"`
	MinProtocolVersion string  `json:"min_protocol_version,omitempty"`
	MaxProtocolVersion string  `json:"max_protocol_version,omitempty"`
}

// PortTrafficPolicy represents port-level traffic policy
type PortTrafficPolicy struct {
	Port             *PortSelector      `json:"port,omitempty"`
	LoadBalancer     *LoadBalancerSettings `json:"load_balancer,omitempty"`
	ConnectionPool   *ConnectionPoolSettings `json:"connection_pool,omitempty"`
	OutlierDetection *OutlierDetection  `json:"outlier_detection,omitempty"`
	TLS              *ClientTLSSettings `json:"tls,omitempty"`
}

// Subset represents a subset definition
type Subset struct {
	Name          string            `json:"name"`
	Labels        map[string]string `json:"labels"`
	TrafficPolicy *TrafficPolicy    `json:"traffic_policy,omitempty"`
}

// Gateway represents Istio Gateway configuration
type Gateway struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Namespace   string           `json:"namespace"`
	Selector    map[string]string `json:"selector"`
	Servers     []Server         `json:"servers"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// Server represents gateway server configuration
type Server struct {
	Port  *Port         `json:"port"`
	Bind  string        `json:"bind,omitempty"`
	Hosts []string      `json:"hosts"`
	TLS   *ServerTLSSettings `json:"tls,omitempty"`
	Name  string        `json:"name,omitempty"`
}

// Port represents port configuration
type Port struct {
	Number   uint32 `json:"number"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"` // HTTP, HTTPS, GRPC, HTTP2, MONGO, TCP, TLS
	TargetPort uint32 `json:"target_port,omitempty"`
}

// ServerTLSSettings represents server TLS configuration
type ServerTLSSettings struct {
	HTTPSRedirect     bool     `json:"https_redirect,omitempty"`
	Mode              string   `json:"mode"` // PASSTHROUGH, SIMPLE, MUTUAL, AUTO_PASSTHROUGH, ISTIO_MUTUAL
	ServerCertificate string   `json:"server_certificate,omitempty"`
	PrivateKey        string   `json:"private_key,omitempty"`
	CACertificates    string   `json:"ca_certificates,omitempty"`
	SubjectAltNames   []string `json:"subject_alt_names,omitempty"`
	CredentialName    string   `json:"credential_name,omitempty"`
	MinProtocolVersion string  `json:"min_protocol_version,omitempty"`
	MaxProtocolVersion string  `json:"max_protocol_version,omitempty"`
	CipherSuites      []string `json:"cipher_suites,omitempty"`
}

// SecurityPolicy represents Istio security policies
type SecurityPolicy struct {
	ID                 uuid.UUID        `json:"id"`
	Name               string           `json:"name"`
	Namespace          string           `json:"namespace"`
	AuthorizationPolicy *AuthorizationPolicy `json:"authorization_policy,omitempty"`
	RequestAuthentication *RequestAuthentication `json:"request_authentication,omitempty"`
	PeerAuthentication *PeerAuthentication `json:"peer_authentication,omitempty"`
	Labels             map[string]string `json:"labels"`
	Annotations        map[string]string `json:"annotations"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

// AuthorizationPolicy represents authorization policy
type AuthorizationPolicy struct {
	Rules    []AuthorizationRule `json:"rules,omitempty"`
	Action   string              `json:"action,omitempty"` // ALLOW, DENY, AUDIT, CUSTOM
	Provider *ExtensionProvider  `json:"provider,omitempty"`
}

// AuthorizationRule represents authorization rule
type AuthorizationRule struct {
	From []AuthorizationRuleFrom `json:"from,omitempty"`
	To   []AuthorizationRuleTo   `json:"to,omitempty"`
	When []AuthorizationRuleWhen `json:"when,omitempty"`
}

// AuthorizationRuleFrom represents source specification
type AuthorizationRuleFrom struct {
	Source *Source `json:"source,omitempty"`
}

// Source represents source specification
type Source struct {
	Principals       []string `json:"principals,omitempty"`
	RequestPrincipals []string `json:"request_principals,omitempty"`
	Namespaces       []string `json:"namespaces,omitempty"`
	IPBlocks         []string `json:"ip_blocks,omitempty"`
	RemoteIPBlocks   []string `json:"remote_ip_blocks,omitempty"`
}

// AuthorizationRuleTo represents destination specification
type AuthorizationRuleTo struct {
	Operation *Operation `json:"operation,omitempty"`
}

// Operation represents operation specification
type Operation struct {
	Methods []string `json:"methods,omitempty"`
	Hosts   []string `json:"hosts,omitempty"`
	Ports   []string `json:"ports,omitempty"`
	Paths   []string `json:"paths,omitempty"`
}

// AuthorizationRuleWhen represents condition specification
type AuthorizationRuleWhen struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// ExtensionProvider represents extension provider
type ExtensionProvider struct {
	Name string `json:"name"`
}

// RequestAuthentication represents request authentication
type RequestAuthentication struct {
	JWTRules []JWTRule `json:"jwt_rules,omitempty"`
}

// JWTRule represents JWT rule
type JWTRule struct {
	Issuer               string   `json:"issuer,omitempty"`
	Audiences            []string `json:"audiences,omitempty"`
	JWKSUri              string   `json:"jwks_uri,omitempty"`
	JWKS                 string   `json:"jwks,omitempty"`
	FromHeaders          []JWTHeader `json:"from_headers,omitempty"`
	FromParams           []string `json:"from_params,omitempty"`
	OutputPayloadToHeader string  `json:"output_payload_to_header,omitempty"`
	ForwardOriginalToken bool    `json:"forward_original_token,omitempty"`
}

// JWTHeader represents JWT header
type JWTHeader struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix,omitempty"`
}

// PeerAuthentication represents peer authentication
type PeerAuthentication struct {
	MTLSMode string                      `json:"mtls_mode"` // UNSET, DISABLE, PERMISSIVE, STRICT
	PortLevelMTLS []PortLevelMTLS        `json:"port_level_mtls,omitempty"`
}

// PortLevelMTLS represents port-level mTLS
type PortLevelMTLS struct {
	Port     *PortSelector `json:"port"`
	MTLSMode string        `json:"mtls_mode"`
}

// Circuit breaker, retry, and timeout configurations
type CircuitBreakerConfig struct {
	MaxConnections     int32         `json:"max_connections"`
	MaxPendingRequests int32         `json:"max_pending_requests"`
	MaxRequests        int32         `json:"max_requests"`
	MaxRetries         int32         `json:"max_retries"`
	ConsecutiveErrors  uint32        `json:"consecutive_errors"`
	Interval           time.Duration `json:"interval"`
	BaseEjectionTime   time.Duration `json:"base_ejection_time"`
	MaxEjectionPercent int32         `json:"max_ejection_percent"`
}

type RetryConfig struct {
	Attempts      int32         `json:"attempts"`
	PerTryTimeout time.Duration `json:"per_try_timeout"`
	RetryOn       string        `json:"retry_on"`
}

type TimeoutConfig struct {
	RequestTimeout  time.Duration `json:"request_timeout"`
	ConnectTimeout  time.Duration `json:"connect_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
}

// NewIstioManager creates a new Istio manager
func NewIstioManager(config *IstioConfig) *IstioManager {
	if config == nil {
		config = &IstioConfig{
			Namespace:        "istio-system",
			IngressGateway:   "istio-ingressgateway",
			TLSMode:          "ISTIO_MUTUAL",
			mTLSEnabled:      true,
			TracingEnabled:   true,
			MetricsEnabled:   true,
			AccessLogEnabled: true,
			CircuitBreaker: &CircuitBreakerConfig{
				MaxConnections:     1000,
				MaxPendingRequests: 100,
				MaxRequests:        100,
				MaxRetries:         3,
				ConsecutiveErrors:  5,
				Interval:           30 * time.Second,
				BaseEjectionTime:   30 * time.Second,
				MaxEjectionPercent: 50,
			},
			Retry: &RetryConfig{
				Attempts:      3,
				PerTryTimeout: 30 * time.Second,
				RetryOn:       "5xx,reset,connect-failure,refused-stream",
			},
			Timeout: &TimeoutConfig{
				RequestTimeout: 60 * time.Second,
				ConnectTimeout: 10 * time.Second,
				IdleTimeout:    60 * time.Second,
			},
		}
	}

	return &IstioManager{
		config:           config,
		virtualServices:  make(map[string]*VirtualService),
		destinationRules: make(map[string]*DestinationRule),
		gateways:         make(map[string]*Gateway),
		policies:         make(map[string]*SecurityPolicy),
	}
}

// CreateVirtualService creates a new virtual service
func (im *IstioManager) CreateVirtualService(vs *VirtualService) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if vs.ID == uuid.Nil {
		vs.ID = uuid.New()
	}

	vs.CreatedAt = time.Now()
	vs.UpdatedAt = time.Now()

	if vs.Namespace == "" {
		vs.Namespace = im.config.Namespace
	}

	im.virtualServices[vs.Name] = vs
	return nil
}

// CreateDestinationRule creates a new destination rule
func (im *IstioManager) CreateDestinationRule(dr *DestinationRule) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if dr.ID == uuid.Nil {
		dr.ID = uuid.New()
	}

	dr.CreatedAt = time.Now()
	dr.UpdatedAt = time.Now()

	if dr.Namespace == "" {
		dr.Namespace = im.config.Namespace
	}

	// Apply default traffic policy if not specified
	if dr.TrafficPolicy == nil {
		dr.TrafficPolicy = im.createDefaultTrafficPolicy()
	}

	im.destinationRules[dr.Name] = dr
	return nil
}

// CreateGateway creates a new gateway
func (im *IstioManager) CreateGateway(gw *Gateway) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if gw.ID == uuid.Nil {
		gw.ID = uuid.New()
	}

	gw.CreatedAt = time.Now()
	gw.UpdatedAt = time.Now()

	if gw.Namespace == "" {
		gw.Namespace = im.config.Namespace
	}

	// Default selector for Istio ingress gateway
	if gw.Selector == nil {
		gw.Selector = map[string]string{"istio": "ingressgateway"}
	}

	im.gateways[gw.Name] = gw
	return nil
}

// CreateSecurityPolicy creates a new security policy
func (im *IstioManager) CreateSecurityPolicy(policy *SecurityPolicy) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if policy.ID == uuid.Nil {
		policy.ID = uuid.New()
	}

	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	if policy.Namespace == "" {
		policy.Namespace = im.config.Namespace
	}

	im.policies[policy.Name] = policy
	return nil
}

// createDefaultTrafficPolicy creates default traffic policy with circuit breaker
func (im *IstioManager) createDefaultTrafficPolicy() *TrafficPolicy {
	return &TrafficPolicy{
		LoadBalancer: &LoadBalancerSettings{
			Simple: "LEAST_CONN",
		},
		ConnectionPool: &ConnectionPoolSettings{
			TCP: &TCPSettings{
				MaxConnections: im.config.CircuitBreaker.MaxConnections,
				ConnectTimeout: im.config.Timeout.ConnectTimeout,
			},
			HTTP: &HTTPSettings{
				HTTP1MaxPendingRequests:  im.config.CircuitBreaker.MaxPendingRequests,
				HTTP2MaxRequests:         im.config.CircuitBreaker.MaxRequests,
				MaxRequestsPerConnection: im.config.CircuitBreaker.MaxRequests,
				MaxRetries:               im.config.CircuitBreaker.MaxRetries,
				IdleTimeout:              im.config.Timeout.IdleTimeout,
			},
		},
		OutlierDetection: &OutlierDetection{
			Consecutive5xxErrors: im.config.CircuitBreaker.ConsecutiveErrors,
			Interval:             im.config.CircuitBreaker.Interval,
			BaseEjectionTime:     im.config.CircuitBreaker.BaseEjectionTime,
			MaxEjectionPercent:   im.config.CircuitBreaker.MaxEjectionPercent,
		},
	}
}

// FastenMind Istio Configurations

// FastenMindGateway creates the main gateway for FastenMind
func FastenMindGateway() *Gateway {
	return &Gateway{
		Name:      "fastenmind-gateway",
		Namespace: "fastenmind",
		Selector:  map[string]string{"istio": "ingressgateway"},
		Servers: []Server{
			{
				Port: &Port{
					Number:   80,
					Name:     "http",
					Protocol: "HTTP",
				},
				Hosts: []string{"api.fastenmind.com", "app.fastenmind.com"},
				TLS: &ServerTLSSettings{
					HTTPSRedirect: true,
				},
			},
			{
				Port: &Port{
					Number:   443,
					Name:     "https",
					Protocol: "HTTPS",
				},
				Hosts: []string{"api.fastenmind.com", "app.fastenmind.com"},
				TLS: &ServerTLSSettings{
					Mode:           "SIMPLE",
					CredentialName: "fastenmind-tls-cert",
				},
			},
		},
		Labels: map[string]string{
			"app":     "fastenmind",
			"version": "v1",
		},
	}
}

// FastenMindAPIVirtualService creates virtual service for API routing
func FastenMindAPIVirtualService() *VirtualService {
	return &VirtualService{
		Name:      "fastenmind-api-vs",
		Namespace: "fastenmind",
		Hosts:     []string{"api.fastenmind.com"},
		Gateways:  []string{"fastenmind-gateway"},
		HTTPRoutes: []HTTPRoute{
			{
				Name: "api-v1",
				Match: []HTTPMatchRequest{
					{
						URI: &StringMatch{Prefix: "/api/v1/"},
					},
				},
				Route: []HTTPRouteDestination{
					{
						Destination: &Destination{
							Host:   "fastenmind-api-service",
							Subset: "v1",
						},
						Weight: 90,
					},
					{
						Destination: &Destination{
							Host:   "fastenmind-api-service",
							Subset: "v2",
						},
						Weight: 10, // Canary deployment
					},
				},
				Timeout: 30 * time.Second,
				Retry: &HTTPRetry{
					Attempts:      3,
					PerTryTimeout: 10 * time.Second,
					RetryOn:       "5xx,reset,connect-failure",
				},
				Headers: &Headers{
					Request: &HeaderOperations{
						Set: map[string]string{
							"X-Service": "fastenmind-api",
						},
					},
				},
			},
			{
				Name: "health-check",
				Match: []HTTPMatchRequest{
					{
						URI: &StringMatch{Exact: "/health"},
					},
				},
				Route: []HTTPRouteDestination{
					{
						Destination: &Destination{
							Host: "fastenmind-api-service",
						},
					},
				},
			},
		},
	}
}

// FastenMindAPIDestinationRule creates destination rule for API service
func FastenMindAPIDestinationRule() *DestinationRule {
	return &DestinationRule{
		Name:      "fastenmind-api-dr",
		Namespace: "fastenmind",
		Host:      "fastenmind-api-service",
		TrafficPolicy: &TrafficPolicy{
			LoadBalancer: &LoadBalancerSettings{
				Simple: "LEAST_CONN",
			},
			ConnectionPool: &ConnectionPoolSettings{
				TCP: &TCPSettings{
					MaxConnections: 100,
					ConnectTimeout: 10 * time.Second,
				},
				HTTP: &HTTPSettings{
					HTTP1MaxPendingRequests:  50,
					HTTP2MaxRequests:         100,
					MaxRequestsPerConnection: 10,
					MaxRetries:               3,
					IdleTimeout:              60 * time.Second,
				},
			},
			OutlierDetection: &OutlierDetection{
				Consecutive5xxErrors: 5,
				Interval:             30 * time.Second,
				BaseEjectionTime:     30 * time.Second,
				MaxEjectionPercent:   50,
			},
		},
		Subsets: []Subset{
			{
				Name:   "v1",
				Labels: map[string]string{"version": "v1"},
			},
			{
				Name:   "v2",
				Labels: map[string]string{"version": "v2"},
				TrafficPolicy: &TrafficPolicy{
					ConnectionPool: &ConnectionPoolSettings{
						HTTP: &HTTPSettings{
							MaxRetries: 1, // More aggressive for canary
						},
					},
				},
			},
		},
	}
}

// FastenMindSecurityPolicy creates security policies for FastenMind
func FastenMindSecurityPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		Name:      "fastenmind-security",
		Namespace: "fastenmind",
		AuthorizationPolicy: &AuthorizationPolicy{
			Action: "ALLOW",
			Rules: []AuthorizationRule{
				{
					From: []AuthorizationRuleFrom{
						{
							Source: &Source{
								Principals: []string{"cluster.local/ns/fastenmind/sa/api-service"},
							},
						},
					},
					To: []AuthorizationRuleTo{
						{
							Operation: &Operation{
								Methods: []string{"GET", "POST", "PUT", "DELETE"},
								Paths:   []string{"/api/v1/*"},
							},
						},
					},
				},
			},
		},
		RequestAuthentication: &RequestAuthentication{
			JWTRules: []JWTRule{
				{
					Issuer:  "https://auth.fastenmind.com",
					JWKSUri: "https://auth.fastenmind.com/.well-known/jwks.json",
					FromHeaders: []JWTHeader{
						{
							Name:   "Authorization",
							Prefix: "Bearer ",
						},
					},
				},
			},
		},
		PeerAuthentication: &PeerAuthentication{
			MTLSMode: "STRICT",
		},
	}
}

// Apply configurations to Kubernetes (placeholder implementation)
func (im *IstioManager) ApplyConfiguration(ctx context.Context) error {
	// This would implement actual Kubernetes API calls to apply Istio configurations
	// Using kubectl or Kubernetes client-go library
	
	// Apply gateways
	for _, gateway := range im.gateways {
		if err := im.applyGateway(ctx, gateway); err != nil {
			return fmt.Errorf("failed to apply gateway %s: %w", gateway.Name, err)
		}
	}

	// Apply virtual services
	for _, vs := range im.virtualServices {
		if err := im.applyVirtualService(ctx, vs); err != nil {
			return fmt.Errorf("failed to apply virtual service %s: %w", vs.Name, err)
		}
	}

	// Apply destination rules
	for _, dr := range im.destinationRules {
		if err := im.applyDestinationRule(ctx, dr); err != nil {
			return fmt.Errorf("failed to apply destination rule %s: %w", dr.Name, err)
		}
	}

	// Apply security policies
	for _, policy := range im.policies {
		if err := im.applySecurityPolicy(ctx, policy); err != nil {
			return fmt.Errorf("failed to apply security policy %s: %w", policy.Name, err)
		}
	}

	return nil
}

// Placeholder implementations for Kubernetes API calls
func (im *IstioManager) applyGateway(ctx context.Context, gateway *Gateway) error {
	// kubectl apply -f gateway.yaml
	return nil
}

func (im *IstioManager) applyVirtualService(ctx context.Context, vs *VirtualService) error {
	// kubectl apply -f virtualservice.yaml
	return nil
}

func (im *IstioManager) applyDestinationRule(ctx context.Context, dr *DestinationRule) error {
	// kubectl apply -f destinationrule.yaml
	return nil
}

func (im *IstioManager) applySecurityPolicy(ctx context.Context, policy *SecurityPolicy) error {
	// kubectl apply -f authorizationpolicy.yaml
	// kubectl apply -f requestauthentication.yaml
	// kubectl apply -f peerauthentication.yaml
	return nil
}