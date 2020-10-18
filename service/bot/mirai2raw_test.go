package bot

import (
	"html"
	"testing"
)

func TestEscape(t *testing.T) {
	result := html.EscapeString("<a />")
	 t.Log(result)
	t.Log(html.UnescapeString(result))
}
