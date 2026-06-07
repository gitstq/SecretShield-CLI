// Package utils provides utility functions for SecretShield-CLI.
package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// IsBinaryFile checks if a file is likely a binary file by reading the first 512 bytes.
func IsBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return true
	}

	for i := 0; i < n; i++ {
		b := buf[i]
		// Check for null byte which indicates binary content
		if b == 0 {
			return true
		}
	}
	return false
}

// IsExcluded checks if a path matches any of the excluded directory patterns.
func IsExcluded(path string, excludes []string) bool {
	if len(excludes) == 0 {
		return false
	}

	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		for _, ex := range excludes {
			if strings.EqualFold(part, ex) {
				return true
			}
		}
	}
	return false
}

// MaskSecret masks a matched secret string for safe display.
// It shows the first 4 and last 4 characters, replacing the middle with asterisks.
func MaskSecret(secret string) string {
	secret = strings.TrimSpace(secret)
	if len(secret) <= 8 {
		return strings.Repeat("*", len(secret))
	}
	return secret[:4] + strings.Repeat("*", len(secret)-8) + secret[len(secret)-4:]
}

// FileExists checks if a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir checks if the path is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ContainsAny checks if a string contains any of the given substrings.
func ContainsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
