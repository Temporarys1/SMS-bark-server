package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

const defaultDataDir = "/data"
const appName = "bark-server"

// GetDataDir returns the appropriate data directory for the current platform
func GetDataDir() string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxDataDir()
	case "darwin":
		return getDarwinDataDir()
	case "windows":
		return getWindowsDataDir()
	default:
		return defaultDataDir
	}
}

// getLinuxDataDir returns the XDG data directory for Linux
func getLinuxDataDir() string {
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, appName)
	}

	return defaultDataDir
}

// getDarwinDataDir returns the Application Support directory for macOS
func getDarwinDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return defaultDataDir
	}

	appSupport := filepath.Join(home, "Library", "Application Support")
	if isDirectoryWriteable(appSupport) {
		return filepath.Join(appSupport, appName)
	}

	return defaultDataDir
}

// getWindowsDataDir returns the AppData directory for Windows
func getWindowsDataDir() string {
	if appData := os.Getenv("APPDATA"); appData != "" {
		if isDirectoryWriteable(appData) {
			return filepath.Join(appData, appName)
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return defaultDataDir
	}

	appData := filepath.Join(home, "AppData", "Roaming")
	if isDirectoryWriteable(appData) {
		return filepath.Join(appData, appName)
	}

	return defaultDataDir
}

// isDirectoryWriteable checks if a directory exists and is writable
func isDirectoryWriteable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return false
	}

	// 尝试创建临时文件来测试可写性
	testFile := filepath.Join(path, ".bark-server-test")
	f, err := os.Create(testFile)
	if err != nil {
		return false
	}
	_ = f.Close()
	_ = os.Remove(testFile)

	return true
}
