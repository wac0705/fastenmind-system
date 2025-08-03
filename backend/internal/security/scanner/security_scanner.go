package scanner

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// SecurityScanner provides SAST/DAST scanning capabilities
type SecurityScanner struct {
	rules       []SecurityRule
	results     []SecurityIssue
	config      *ScanConfig
	fileFilters []string
}

// ScanConfig holds scanner configuration
type ScanConfig struct {
	// Scan types
	EnableSAST bool
	EnableDAST bool
	EnableDeps bool
	
	// Paths
	SourcePath   string
	ExcludePaths []string
	
	// Severity levels to report
	MinSeverity Severity
	
	// Output
	OutputFormat string // "json", "sarif", "markdown"
	OutputFile   string
}

// SecurityRule represents a security scanning rule
type SecurityRule struct {
	ID          string
	Name        string
	Description string
	Severity    Severity
	Category    string
	Pattern     *regexp.Regexp
	FileTypes   []string
	Check       func(file string, content string) []SecurityIssue
}

// SecurityIssue represents a security finding
type SecurityIssue struct {
	RuleID      string    `json:"rule_id"`
	Severity    Severity  `json:"severity"`
	Category    string    `json:"category"`
	Message     string    `json:"message"`
	File        string    `json:"file"`
	Line        int       `json:"line"`
	Column      int       `json:"column"`
	Code        string    `json:"code"`
	Suggestion  string    `json:"suggestion"`
	Confidence  float64   `json:"confidence"`
	Timestamp   time.Time `json:"timestamp"`
}

// Severity levels
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// NewSecurityScanner creates a new security scanner
func NewSecurityScanner(config *ScanConfig) *SecurityScanner {
	scanner := &SecurityScanner{
		config: config,
		fileFilters: []string{
			"*.go", "*.js", "*.ts", "*.jsx", "*.tsx",
			"*.java", "*.py", "*.rb", "*.php", "*.cs",
			"*.sql", "*.yaml", "*.yml", "*.json", "*.xml",
		},
	}
	
	scanner.initializeRules()
	return scanner
}

// initializeRules sets up security scanning rules
func (s *SecurityScanner) initializeRules() {
	s.rules = []SecurityRule{
		// SQL Injection
		{
			ID:          "SQL001",
			Name:        "SQL Injection",
			Description: "Potential SQL injection vulnerability",
			Severity:    SeverityCritical,
			Category:    "Injection",
			Pattern:     regexp.MustCompile(`(?i)(exec|execute|query|prepare)\s*\(\s*[^"]*\+|fmt\.Sprintf.*select|insert|update|delete`),
			FileTypes:   []string{".go", ".java", ".cs", ".php"},
			Check:       s.checkSQLInjection,
		},
		// Hard-coded Secrets
		{
			ID:          "SEC001",
			Name:        "Hardcoded Secret",
			Description: "Potential hardcoded secret or API key",
			Severity:    SeverityHigh,
			Category:    "Secrets",
			Pattern:     regexp.MustCompile(`(?i)(api[_-]?key|secret|password|token|private[_-]?key)\s*[:=]\s*["'][^"']+["']`),
			FileTypes:   []string{".go", ".js", ".py", ".java", ".env"},
			Check:       s.checkHardcodedSecrets,
		},
		// Weak Cryptography
		{
			ID:          "CRYPTO001",
			Name:        "Weak Cryptography",
			Description: "Use of weak cryptographic algorithm",
			Severity:    SeverityHigh,
			Category:    "Cryptography",
			Pattern:     regexp.MustCompile(`(?i)(md5|sha1|des|rc4)\.`),
			FileTypes:   []string{".go", ".java", ".cs", ".py"},
			Check:       s.checkWeakCrypto,
		},
		// Command Injection
		{
			ID:          "CMD001",
			Name:        "Command Injection",
			Description: "Potential command injection vulnerability",
			Severity:    SeverityCritical,
			Category:    "Injection",
			Pattern:     regexp.MustCompile(`(?i)(exec\.command|os\.system|subprocess\.call)\s*\([^)]*\+`),
			FileTypes:   []string{".go", ".py", ".rb", ".php"},
			Check:       s.checkCommandInjection,
		},
		// Path Traversal
		{
			ID:          "PATH001",
			Name:        "Path Traversal",
			Description: "Potential path traversal vulnerability",
			Severity:    SeverityHigh,
			Category:    "File Access",
			Pattern:     regexp.MustCompile(`(?i)(\.\.\/|\.\.\\\\|filepath\.join.*\+)`),
			FileTypes:   []string{".go", ".java", ".cs", ".py"},
			Check:       s.checkPathTraversal,
		},
		// Insecure Random
		{
			ID:          "RAND001",
			Name:        "Insecure Random",
			Description: "Use of insecure random number generator",
			Severity:    SeverityMedium,
			Category:    "Cryptography",
			Pattern:     regexp.MustCompile(`(?i)(math/rand|random\.random|rand\(\))`),
			FileTypes:   []string{".go", ".py", ".java"},
			Check:       s.checkInsecureRandom,
		},
		// XXE
		{
			ID:          "XXE001",
			Name:        "XML External Entity",
			Description: "Potential XXE vulnerability",
			Severity:    SeverityHigh,
			Category:    "Injection",
			Pattern:     regexp.MustCompile(`(?i)(xml.*external|dtd.*system|entity.*system)`),
			FileTypes:   []string{".go", ".java", ".cs", ".py"},
			Check:       s.checkXXE,
		},
		// Open Redirect
		{
			ID:          "REDIRECT001",
			Name:        "Open Redirect",
			Description: "Potential open redirect vulnerability",
			Severity:    SeverityMedium,
			Category:    "Access Control",
			Pattern:     regexp.MustCompile(`(?i)(redirect|location\.href)\s*=.*request\.|http\.redirect.*request\.`),
			FileTypes:   []string{".go", ".js", ".java", ".php"},
			Check:       s.checkOpenRedirect,
		},
	}
}

// ScanProject scans the entire project
func (s *SecurityScanner) ScanProject(ctx context.Context) error {
	s.results = []SecurityIssue{}
	
	err := filepath.Walk(s.config.SourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Check if cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		// Skip excluded paths
		for _, exclude := range s.config.ExcludePaths {
			if strings.Contains(path, exclude) {
				return nil
			}
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Check file extension
		if !s.shouldScanFile(path) {
			return nil
		}
		
		// Scan file
		if err := s.scanFile(path); err != nil {
			fmt.Printf("Error scanning %s: %v\n", path, err)
		}
		
		return nil
	})
	
	return err
}

// scanFile scans a single file
func (s *SecurityScanner) scanFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	contentStr := string(content)
	
	// Apply rules
	for _, rule := range s.rules {
		// Check file type
		if !s.matchesFileType(filename, rule.FileTypes) {
			continue
		}
		
		// Run custom check if available
		if rule.Check != nil {
			issues := rule.Check(filename, contentStr)
			s.results = append(s.results, issues...)
		} else if rule.Pattern != nil {
			// Use pattern matching
			matches := rule.Pattern.FindAllStringIndex(contentStr, -1)
			for _, match := range matches {
				line, col := s.getLineColumn(contentStr, match[0])
				issue := SecurityIssue{
					RuleID:     rule.ID,
					Severity:   rule.Severity,
					Category:   rule.Category,
					Message:    rule.Description,
					File:       filename,
					Line:       line,
					Column:     col,
					Code:       s.getCodeSnippet(contentStr, line),
					Confidence: 0.8,
					Timestamp:  time.Now(),
				}
				s.results = append(s.results, issue)
			}
		}
	}
	
	// Special checks for Go files
	if strings.HasSuffix(filename, ".go") {
		s.scanGoFile(filename, contentStr)
	}
	
	return nil
}

// scanGoFile performs Go-specific security checks
func (s *SecurityScanner) scanGoFile(filename, content string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return
	}
	
	// Walk AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			s.checkGoCallExpr(fset, x)
		case *ast.AssignStmt:
			s.checkGoAssignment(fset, x)
		}
		return true
	})
}

// Security check implementations

func (s *SecurityScanner) checkSQLInjection(file, content string) []SecurityIssue {
	var issues []SecurityIssue
	
	// Look for string concatenation in SQL queries
	sqlPattern := regexp.MustCompile(`(?i)(query|exec|prepare)\s*\(\s*"[^"]*"\s*\+`)
	matches := sqlPattern.FindAllStringIndex(content, -1)
	
	for _, match := range matches {
		line, col := s.getLineColumn(content, match[0])
		issues = append(issues, SecurityIssue{
			RuleID:     "SQL001",
			Severity:   SeverityCritical,
			Category:   "Injection",
			Message:    "Potential SQL injection - avoid string concatenation in queries",
			File:       file,
			Line:       line,
			Column:     col,
			Code:       s.getCodeSnippet(content, line),
			Suggestion: "Use parameterized queries or prepared statements",
			Confidence: 0.9,
			Timestamp:  time.Now(),
		})
	}
	
	return issues
}

func (s *SecurityScanner) checkHardcodedSecrets(file, content string) []SecurityIssue {
	var issues []SecurityIssue
	
	// Skip test files
	if strings.Contains(file, "_test.go") || strings.Contains(file, "test_") {
		return issues
	}
	
	// Common patterns for secrets
	patterns := []struct {
		pattern *regexp.Regexp
		message string
	}{
		{
			pattern: regexp.MustCompile(`(?i)api[_-]?key\s*[:=]\s*["'][a-zA-Z0-9]{20,}["']`),
			message: "Hardcoded API key detected",
		},
		{
			pattern: regexp.MustCompile(`(?i)password\s*[:=]\s*["'][^"']{8,}["']`),
			message: "Hardcoded password detected",
		},
		{
			pattern: regexp.MustCompile(`(?i)secret\s*[:=]\s*["'][a-zA-Z0-9]{10,}["']`),
			message: "Hardcoded secret detected",
		},
		{
			pattern: regexp.MustCompile(`-----BEGIN (RSA |EC )?PRIVATE KEY-----`),
			message: "Private key detected in source code",
		},
	}
	
	for _, p := range patterns {
		matches := p.pattern.FindAllStringIndex(content, -1)
		for _, match := range matches {
			line, col := s.getLineColumn(content, match[0])
			issues = append(issues, SecurityIssue{
				RuleID:     "SEC001",
				Severity:   SeverityHigh,
				Category:   "Secrets",
				Message:    p.message,
				File:       file,
				Line:       line,
				Column:     col,
				Code:       s.getCodeSnippet(content, line),
				Suggestion: "Use environment variables or secure key management service",
				Confidence: 0.95,
				Timestamp:  time.Now(),
			})
		}
	}
	
	return issues
}

func (s *SecurityScanner) checkWeakCrypto(file, content string) []SecurityIssue {
	var issues []SecurityIssue
	
	weakAlgos := map[string]string{
		"md5":    "MD5 is cryptographically broken",
		"sha1":   "SHA1 is deprecated for security use",
		"des":    "DES is weak encryption",
		"rc4":    "RC4 has known vulnerabilities",
		"ecb":    "ECB mode is insecure",
	}
	
	for algo, message := range weakAlgos {
		pattern := regexp.MustCompile(`(?i)\b` + algo + `\b`)
		matches := pattern.FindAllStringIndex(content, -1)
		for _, match := range matches {
			line, col := s.getLineColumn(content, match[0])
			issues = append(issues, SecurityIssue{
				RuleID:     "CRYPTO001",
				Severity:   SeverityHigh,
				Category:   "Cryptography",
				Message:    message,
				File:       file,
				Line:       line,
				Column:     col,
				Code:       s.getCodeSnippet(content, line),
				Suggestion: "Use SHA256, SHA3, or AES-GCM for security purposes",
				Confidence: 0.85,
				Timestamp:  time.Now(),
			})
		}
	}
	
	return issues
}

func (s *SecurityScanner) checkCommandInjection(file, content string) []SecurityIssue {
	var issues []SecurityIssue
	
	// Check for command execution with user input
	cmdPattern := regexp.MustCompile(`(?i)(exec\.command|os\.system|subprocess\.)\s*\([^)]*\+[^)]*\)`)
	matches := cmdPattern.FindAllStringIndex(content, -1)
	
	for _, match := range matches {
		line, col := s.getLineColumn(content, match[0])
		issues = append(issues, SecurityIssue{
			RuleID:     "CMD001",
			Severity:   SeverityCritical,
			Category:   "Injection",
			Message:    "Potential command injection vulnerability",
			File:       file,
			Line:       line,
			Column:     col,
			Code:       s.getCodeSnippet(content, line),
			Suggestion: "Validate and sanitize user input, use allowlists for commands",
			Confidence: 0.8,
			Timestamp:  time.Now(),
		})
	}
	
	return issues
}

func (s *SecurityScanner) checkPathTraversal(file, content string) []SecurityIssue {
	var issues []SecurityIssue
	
	// Check for path traversal patterns
	traversalPattern := regexp.MustCompile(`(\.\.\/|\.\.\\\\|filepath\.Join\([^)]*\+)`)
	matches := traversalPattern.FindAllStringIndex(content, -1)
	
	for _, match := range matches {
		line, col := s.getLineColumn(content, match[0])
		issues = append(issues, SecurityIssue{
			RuleID:     "PATH001",
			Severity:   SeverityHigh,
			Category:   "File Access",
			Message:    "Potential path traversal vulnerability",
			File:       file,
			Line:       line,
			Column:     col,
			Code:       s.getCodeSnippet(content, line),
			Suggestion: "Use filepath.Clean() and validate paths against allowed directories",
			Confidence: 0.75,
			Timestamp:  time.Now(),
		})
	}
	
	return issues
}

func (s *SecurityScanner) checkInsecureRandom(file, content string) []SecurityIssue {
	var issues []SecurityIssue
	
	// Skip if it's clearly not for security purposes
	if strings.Contains(content, "// not for security") || strings.Contains(content, "// testing only") {
		return issues
	}
	
	insecurePattern := regexp.MustCompile(`(?i)(math/rand|rand\.new|random\.random)`)
	matches := insecurePattern.FindAllStringIndex(content, -1)
	
	for _, match := range matches {
		line, col := s.getLineColumn(content, match[0])
		
		// Check if it's near security-related code
		contextStart := match[0] - 200
		if contextStart < 0 {
			contextStart = 0
		}
		contextEnd := match[1] + 200
		if contextEnd > len(content) {
			contextEnd = len(content)
		}
		context := content[contextStart:contextEnd]
		
		if strings.Contains(strings.ToLower(context), "token") || 
		   strings.Contains(strings.ToLower(context), "key") ||
		   strings.Contains(strings.ToLower(context), "password") {
			issues = append(issues, SecurityIssue{
				RuleID:     "RAND001",
				Severity:   SeverityMedium,
				Category:   "Cryptography",
				Message:    "Insecure random number generator used in security context",
				File:       file,
				Line:       line,
				Column:     col,
				Code:       s.getCodeSnippet(content, line),
				Suggestion: "Use crypto/rand for security-sensitive random values",
				Confidence: 0.9,
				Timestamp:  time.Now(),
			})
		}
	}
	
	return issues
}

func (s *SecurityScanner) checkXXE(file, content string) []SecurityIssue {
	var issues []SecurityIssue
	
	xxePatterns := []struct {
		pattern *regexp.Regexp
		message string
	}{
		{
			pattern: regexp.MustCompile(`(?i)xml.*externalentity.*true`),
			message: "XML external entities are enabled",
		},
		{
			pattern: regexp.MustCompile(`(?i)setfeature.*external.*true`),
			message: "External DTD processing is enabled",
		},
	}
	
	for _, p := range xxePatterns {
		matches := p.pattern.FindAllStringIndex(content, -1)
		for _, match := range matches {
			line, col := s.getLineColumn(content, match[0])
			issues = append(issues, SecurityIssue{
				RuleID:     "XXE001",
				Severity:   SeverityHigh,
				Category:   "Injection",
				Message:    p.message,
				File:       file,
				Line:       line,
				Column:     col,
				Code:       s.getCodeSnippet(content, line),
				Suggestion: "Disable XML external entity processing",
				Confidence: 0.9,
				Timestamp:  time.Now(),
			})
		}
	}
	
	return issues
}

func (s *SecurityScanner) checkOpenRedirect(file, content string) []SecurityIssue {
	var issues []SecurityIssue
	
	redirectPattern := regexp.MustCompile(`(?i)(redirect|location\.href)\s*[=:]\s*[^"']*request\.(query|form|param)`)
	matches := redirectPattern.FindAllStringIndex(content, -1)
	
	for _, match := range matches {
		line, col := s.getLineColumn(content, match[0])
		issues = append(issues, SecurityIssue{
			RuleID:     "REDIRECT001",
			Severity:   SeverityMedium,
			Category:   "Access Control",
			Message:    "Potential open redirect vulnerability",
			File:       file,
			Line:       line,
			Column:     col,
			Code:       s.getCodeSnippet(content, line),
			Suggestion: "Validate redirect URLs against an allowlist",
			Confidence: 0.7,
			Timestamp:  time.Now(),
		})
	}
	
	return issues
}

// Helper methods

func (s *SecurityScanner) shouldScanFile(filename string) bool {
	ext := filepath.Ext(filename)
	for _, filter := range s.fileFilters {
		if strings.HasSuffix(filter, ext) {
			return true
		}
	}
	return false
}

func (s *SecurityScanner) matchesFileType(filename string, fileTypes []string) bool {
	for _, ft := range fileTypes {
		if strings.HasSuffix(filename, ft) {
			return true
		}
	}
	return false
}

func (s *SecurityScanner) getLineColumn(content string, offset int) (line, column int) {
	line = 1
	column = 1
	
	for i := 0; i < offset && i < len(content); i++ {
		if content[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	
	return line, column
}

func (s *SecurityScanner) getCodeSnippet(content string, line int) string {
	lines := strings.Split(content, "\n")
	if line <= 0 || line > len(lines) {
		return ""
	}
	
	start := line - 3
	if start < 0 {
		start = 0
	}
	
	end := line + 2
	if end > len(lines) {
		end = len(lines)
	}
	
	snippet := strings.Join(lines[start:end], "\n")
	return strings.TrimSpace(snippet)
}

func (s *SecurityScanner) checkGoCallExpr(fset *token.FileSet, call *ast.CallExpr) {
	// Check for dangerous function calls
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			funcName := ident.Name + "." + sel.Sel.Name
			
			dangerousFuncs := map[string]string{
				"exec.Command":    "Command injection risk",
				"os.Remove":       "File deletion risk",
				"sql.Query":       "SQL injection risk",
				"template.HTML":   "XSS risk",
			}
			
			if msg, isDangerous := dangerousFuncs[funcName]; isDangerous {
				pos := fset.Position(call.Pos())
				s.results = append(s.results, SecurityIssue{
					RuleID:     "GOFUNC001",
					Severity:   SeverityHigh,
					Category:   "Dangerous Function",
					Message:    msg + " with " + funcName,
					File:       pos.Filename,
					Line:       pos.Line,
					Column:     pos.Column,
					Confidence: 0.7,
					Timestamp:  time.Now(),
				})
			}
		}
	}
}

func (s *SecurityScanner) checkGoAssignment(fset *token.FileSet, assign *ast.AssignStmt) {
	// Check for error handling
	if len(assign.Rhs) == 1 {
		if call, ok := assign.Rhs[0].(*ast.CallExpr); ok {
			// Check if function returns error
			if len(assign.Lhs) >= 2 {
				if ident, ok := assign.Lhs[len(assign.Lhs)-1].(*ast.Ident); ok {
					if ident.Name == "_" {
						pos := fset.Position(assign.Pos())
						s.results = append(s.results, SecurityIssue{
							RuleID:     "GOERR001",
							Severity:   SeverityMedium,
							Category:   "Error Handling",
							Message:    "Ignored error return value",
							File:       pos.Filename,
							Line:       pos.Line,
							Column:     pos.Column,
							Suggestion: "Handle errors appropriately",
							Confidence: 0.9,
							Timestamp:  time.Now(),
						})
					}
				}
			}
		}
	}
}

// GenerateReport generates a security scan report
func (s *SecurityScanner) GenerateReport() (string, error) {
	switch s.config.OutputFormat {
	case "json":
		return s.generateJSONReport()
	case "sarif":
		return s.generateSARIFReport()
	case "markdown":
		return s.generateMarkdownReport()
	default:
		return s.generateMarkdownReport()
	}
}

func (s *SecurityScanner) generateJSONReport() (string, error) {
	report := map[string]interface{}{
		"scan_time": time.Now(),
		"config":    s.config,
		"summary": map[string]int{
			"total":    len(s.results),
			"critical": s.countBySeverity(SeverityCritical),
			"high":     s.countBySeverity(SeverityHigh),
			"medium":   s.countBySeverity(SeverityMedium),
			"low":      s.countBySeverity(SeverityLow),
			"info":     s.countBySeverity(SeverityInfo),
		},
		"issues": s.results,
	}
	
	data, err := json.MarshalIndent(report, "", "  ")
	return string(data), err
}

func (s *SecurityScanner) generateSARIFReport() (string, error) {
	// SARIF 2.1.0 format
	sarif := map[string]interface{}{
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"version": "2.1.0",
		"runs": []interface{}{
			map[string]interface{}{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":    "FastenMind Security Scanner",
						"version": "1.0.0",
						"rules":   s.rulesToSARIF(),
					},
				},
				"results": s.resultsToSARIF(),
			},
		},
	}
	
	data, err := json.MarshalIndent(sarif, "", "  ")
	return string(data), err
}

func (s *SecurityScanner) generateMarkdownReport() (string, error) {
	var report strings.Builder
	
	report.WriteString("# Security Scan Report\n\n")
	report.WriteString(fmt.Sprintf("**Scan Date:** %s\n\n", time.Now().Format(time.RFC3339)))
	
	// Summary
	report.WriteString("## Summary\n\n")
	report.WriteString(fmt.Sprintf("- **Total Issues:** %d\n", len(s.results)))
	report.WriteString(fmt.Sprintf("- **Critical:** %d\n", s.countBySeverity(SeverityCritical)))
	report.WriteString(fmt.Sprintf("- **High:** %d\n", s.countBySeverity(SeverityHigh)))
	report.WriteString(fmt.Sprintf("- **Medium:** %d\n", s.countBySeverity(SeverityMedium)))
	report.WriteString(fmt.Sprintf("- **Low:** %d\n", s.countBySeverity(SeverityLow)))
	report.WriteString(fmt.Sprintf("- **Info:** %d\n\n", s.countBySeverity(SeverityInfo)))
	
	// Issues by severity
	severities := []Severity{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo}
	severityNames := []string{"Critical", "High", "Medium", "Low", "Info"}
	
	for i, severity := range severities {
		issues := s.filterBySeverity(severity)
		if len(issues) > 0 {
			report.WriteString(fmt.Sprintf("## %s Severity Issues\n\n", severityNames[i]))
			
			for _, issue := range issues {
				report.WriteString(fmt.Sprintf("### %s - %s\n\n", issue.RuleID, issue.Message))
				report.WriteString(fmt.Sprintf("- **File:** %s:%d:%d\n", issue.File, issue.Line, issue.Column))
				report.WriteString(fmt.Sprintf("- **Category:** %s\n", issue.Category))
				report.WriteString(fmt.Sprintf("- **Confidence:** %.0f%%\n", issue.Confidence*100))
				
				if issue.Code != "" {
					report.WriteString("\n**Code:**\n```\n")
					report.WriteString(issue.Code)
					report.WriteString("\n```\n")
				}
				
				if issue.Suggestion != "" {
					report.WriteString(fmt.Sprintf("\n**Suggestion:** %s\n", issue.Suggestion))
				}
				
				report.WriteString("\n---\n\n")
			}
		}
	}
	
	return report.String(), nil
}

func (s *SecurityScanner) countBySeverity(severity Severity) int {
	count := 0
	for _, issue := range s.results {
		if issue.Severity == severity {
			count++
		}
	}
	return count
}

func (s *SecurityScanner) filterBySeverity(severity Severity) []SecurityIssue {
	var filtered []SecurityIssue
	for _, issue := range s.results {
		if issue.Severity == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func (s *SecurityScanner) rulesToSARIF() []interface{} {
	var rules []interface{}
	for _, rule := range s.rules {
		rules = append(rules, map[string]interface{}{
			"id":               rule.ID,
			"name":             rule.Name,
			"shortDescription": map[string]string{"text": rule.Description},
			"defaultConfiguration": map[string]interface{}{
				"level": s.severityToSARIFLevel(rule.Severity),
			},
		})
	}
	return rules
}

func (s *SecurityScanner) resultsToSARIF() []interface{} {
	var results []interface{}
	for _, issue := range s.results {
		results = append(results, map[string]interface{}{
			"ruleId": issue.RuleID,
			"level":  s.severityToSARIFLevel(issue.Severity),
			"message": map[string]string{
				"text": issue.Message,
			},
			"locations": []interface{}{
				map[string]interface{}{
					"physicalLocation": map[string]interface{}{
						"artifactLocation": map[string]string{
							"uri": issue.File,
						},
						"region": map[string]interface{}{
							"startLine":   issue.Line,
							"startColumn": issue.Column,
						},
					},
				},
			},
		})
	}
	return results
}

func (s *SecurityScanner) severityToSARIFLevel(severity Severity) string {
	switch severity {
	case SeverityCritical, SeverityHigh:
		return "error"
	case SeverityMedium:
		return "warning"
	case SeverityLow, SeverityInfo:
		return "note"
	default:
		return "note"
	}
}