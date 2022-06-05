package links

import (
	"testing"
)

type entry struct {
	srcAddress string
	srcDisplay string
	srcHypha   string

	kind    LinkType
	href    string
	display string
}

func TestLink(t *testing.T) {
	// address — display — srchypha
	mappings := []entry{
		{"apple", "", "home",
			LinkLocalHypha, "/hypha/apple", "apple"},
		{"Apple", "Яблоко", "home",
			LinkLocalHypha, "/hypha/apple", "Яблоко"},
		{"   apple ", "Pomme", "terra/incognita",
			LinkLocalHypha, "/hypha/apple", "Pomme"},
		{"./meme", "", "apple",
			LinkLocalHypha, "/hypha/apple/meme", "./meme"},
		{"https://app.le", "app.le website", "app_dot_le",
			LinkExternal, "https://app.le", "app.le website"},
	}
	for _, mapping := range mappings {
		var (
			link    = LegacyFrom(mapping.srcAddress, mapping.srcDisplay, mapping.srcHypha)
			kind    = link.kind
			href    = link.Href()
			display = link.Display()
		)
		if kind != mapping.kind {
			t.Errorf(`When parsing %s→%s@%s got wrong kind: %v`, mapping.srcAddress, mapping.srcDisplay, mapping.srcHypha, kind)
		}
		if href != mapping.href {
			t.Errorf(`When parsing %s→%s@%s got wrong href: [%s], expected [%s]`, mapping.srcAddress, mapping.srcDisplay, mapping.srcHypha, href, mapping.href)
		}
		if display != mapping.display {
			t.Errorf(`When parsing %s→%s@%s got wrong display: [%s], expected [%s]`, mapping.srcAddress, mapping.srcDisplay, mapping.srcHypha, display, mapping.display)
		}
	}
}
