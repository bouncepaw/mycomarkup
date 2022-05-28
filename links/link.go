package links

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v4/mycocontext"
	"path"
	"strings"
)

// Link is a link of some kind.
type Link interface {
	// Classes returns a string to put into the class attr in HTML.
	Classes(ctx mycocontext.Context) string

	// LinkHref returns a string to put into the href attr of <a>.
	LinkHref(ctx mycocontext.Context) string

	// ImgSrc returns a string to put into the src attr of <img>.
	ImgSrc(ctx mycocontext.Context) string

	// DisplayedText returns a string to put inside <a>.
	DisplayedText() string

	// HyphaProbe returns a function that captures the Link. Probes are checked against all existing hyphae. This is Mycorrhiza-specific. If it is nil, do not check this link for existence. TODO: make it optional.
	HyphaProbe() func(string)
}

func LinkFrom(ctx mycocontext.Context, target, display string) Link {
	target, display = strings.TrimSpace(target), strings.TrimSpace(display)
	switch {
	case strings.ContainsRune(target, ':'):
		return &URLLink{
			target:  target,
			display: display,
		}
	case strings.HasPrefix(target, "/"):
		return &LocalRootedLink{
			target:  target,
			display: display,
		}
	case strings.ContainsRune(target, '>'):
		gtpos := strings.IndexRune(target, '>')
		return &InterwikiLink{
			prefix:  target[:gtpos],
			target:  target[gtpos+1:],
			display: display,
		}
	case target == "..":
		return &LocalLink{
			target:  path.Dir(ctx.HyphaName()),
			display: display,
		}
	case strings.HasPrefix(target, "./"):
		var anchor string
		if hashPos := strings.IndexRune(target, '#'); hashPos != -1 {
			anchor = target[hashPos+1:]
			target = target[:hashPos]
		}
		return &LocalLink{
			target:  path.Join(ctx.HyphaName(), target[2:]),
			display: display,
			anchor:  anchor,
		}
	case strings.HasPrefix(target, "../"):
		var anchor string
		if hashPos := strings.IndexRune(target, '#'); hashPos != -1 {
			anchor = target[hashPos+1:]
			target = target[:hashPos]
		}
		return &LocalLink{
			target:  path.Join(path.Dir(ctx.HyphaName()), target[3:]),
			display: display,
			anchor:  anchor,
		}
	case strings.ContainsRune(target, '#'):
		hashPos := strings.IndexRune(target, '#')
		anchor := target[hashPos+1:]
		target = target[:hashPos]
		return &LocalLink{
			target:  target,
			display: display,
			anchor:  anchor,
		}
	default:
		return &LocalLink{
			target:  target,
			display: display,
		}
	}
}

type LocalLink struct {
	target   string
	display  string
	anchor   string
	existing bool
}

func (l *LocalLink) Classes(ctx mycocontext.Context) string {
	res := "wikilink wikilink_internal"
	if !l.existing && mycocontext.Options(ctx).RedLinksSupported {
		res += " wikilink_new"
	}
	return res
}

func (l *LocalLink) LinkHref(ctx mycocontext.Context) string {
	if l.anchor != "" {
		return mycocontext.Options(ctx).LocalLinkHref(l.target + "#" + l.anchor)
	}
	return mycocontext.Options(ctx).LocalLinkHref(l.target)
}

func (l *LocalLink) ImgSrc(ctx mycocontext.Context) string {
	if l.anchor != "" {
		return mycocontext.Options(ctx).LocalImgSrc(l.target + "#" + l.anchor)
	}
	return mycocontext.Options(ctx).LocalImgSrc(l.target)
}

func (l *LocalLink) DisplayedText() string {
	return l.display
}

func (l *LocalLink) HyphaProbe() func(string) {
	done := false
	return func(docName string) {
		if done {
			return
		}
		if docName == l.target {
			l.existing = true
			done = true
		}
	}
}

type LocalRootedLink struct {
	target, display string
}

func (l *LocalRootedLink) Classes(ctx mycocontext.Context) string {
	return "wikilink wikilink_internal"
}

func (l *LocalRootedLink) LinkHref(ctx mycocontext.Context) string {
	return l.target
}

func (l *LocalRootedLink) ImgSrc(ctx mycocontext.Context) string {
	return l.target
}

func (l *LocalRootedLink) DisplayedText() string {
	return l.display
}

func (l *LocalRootedLink) HyphaProbe() func(string) {
	return nil
}

type URLLink struct {
	target  string
	display string
}

func (U *URLLink) protocol() string {
	return U.target[:strings.IndexRune(U.target, ':')]
}

func (U *URLLink) Classes(ctx mycocontext.Context) string {
	return fmt.Sprintf(
		"wikilink wikilink_external wikilink_%s",
		U.protocol(),
	)
}

func (U *URLLink) LinkHref(ctx mycocontext.Context) string {
	return U.target
}

func (U *URLLink) ImgSrc(ctx mycocontext.Context) string {
	return U.target
}

func (U *URLLink) DisplayedText() string {
	return U.display
}

func (U *URLLink) HyphaProbe() func(string) {
	return nil
}

/*
InterwikiLink in Mycomarkup has this syntax:

	[[prefix>target]]
	[[prefix>target|display]]

For every prefix, there is a known link format. A link format is a format string, that might resemble Go's format strings, but they are actually not. This is DSL for link formats. It is inspired by DokuWiki's interwiki link format: https://www.dokuwiki.org/interwiki.

	https://example.org/view/{NAME}

Supported instructions are (more will be added):

	{NAME} is the document name without any encoding.
*/
type InterwikiLink struct {
	prefix, target, display string
}

func (l *InterwikiLink) Classes(ctx mycocontext.Context) string {
	return "wikilink wikilink_interwiki"
}

func (l *InterwikiLink) LinkHref(ctx mycocontext.Context) string {
	format := mycocontext.Options(ctx).LinkHrefFormatForInterwikiPrefix(l.prefix)
	return strings.ReplaceAll(format, "{NAME}", l.target)
}

func (l *InterwikiLink) ImgSrc(ctx mycocontext.Context) string {
	format := mycocontext.Options(ctx).ImgSrcFormatForInterwikiPrefix(l.prefix)
	return strings.ReplaceAll(format, "{NAME}", l.target)
}

func (l *InterwikiLink) DisplayedText() string {
	return l.display
}

func (l *InterwikiLink) HyphaProbe() func(string) {
	return nil
}
