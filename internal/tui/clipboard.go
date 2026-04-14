package tui

import "github.com/atotto/clipboard"

func copyClipboard(s string) error {
	return clipboard.WriteAll(s)
}
