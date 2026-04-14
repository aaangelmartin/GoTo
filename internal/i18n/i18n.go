// Package i18n provides a minimal bilingual (English/Spanish) string catalog.
//
// The active language is resolved once, in this order:
//  1. The --lang flag (wired from the CLI via SetLang).
//  2. The GOTO_LANG environment variable (values: "en" or "es").
//  3. The LANG / LC_ALL environment variable (anything starting with "es" is Spanish).
//  4. Fallback: English.
//
// Paquete i18n: catálogo mínimo de cadenas bilingüe (inglés/español).
//
// El idioma activo se resuelve una sola vez, en este orden:
//  1. Bandera --lang (conectada por la CLI vía SetLang).
//  2. Variable de entorno GOTO_LANG ("en" o "es").
//  3. Variables LANG / LC_ALL (cualquiera que empiece por "es" es español).
//  4. Fallback: inglés.
package i18n

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Lang is a supported language code.
type Lang string

const (
	// EN is English.
	EN Lang = "en"
	// ES is Spanish.
	ES Lang = "es"
)

var (
	mu      sync.RWMutex
	current Lang
	once    sync.Once
)

// Current returns the active language, resolving it lazily on first call.
func Current() Lang {
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		current = detect()
	})
	mu.RLock()
	defer mu.RUnlock()
	return current
}

// SetLang overrides the active language. Use "" to reset to auto-detection.
func SetLang(l string) {
	mu.Lock()
	defer mu.Unlock()
	switch strings.ToLower(strings.TrimSpace(l)) {
	case "es", "spanish", "español", "espanol":
		current = ES
	case "en", "english":
		current = EN
	case "":
		current = detect()
	default:
		current = EN
	}
	// Mark resolved so Current doesn't overwrite.
	once.Do(func() {})
}

func detect() Lang {
	if v := os.Getenv("GOTO_LANG"); v != "" {
		if strings.HasPrefix(strings.ToLower(v), "es") {
			return ES
		}
		return EN
	}
	for _, k := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if v := os.Getenv(k); v != "" {
			if strings.HasPrefix(strings.ToLower(v), "es") {
				return ES
			}
			return EN
		}
	}
	return EN
}

// T returns the translation for the current language, falling back to English.
func T(key string) string {
	if s, ok := catalog[Current()][key]; ok {
		return s
	}
	if s, ok := catalog[EN][key]; ok {
		return s
	}
	return key
}

// Tf is a convenience wrapper that applies sprintf-style formatting to T(key).
// Uso: i18n.Tf("added_alias", name, url)
func Tf(key string, args ...any) string {
	return fmt.Sprintf(T(key), args...)
}
