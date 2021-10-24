package genhtml

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/genhtml/tag"
	"github.com/bouncepaw/mycomarkup/v3/util/lines"
	"log"
)

func wrapInTransclusionError(errParagraph string) tag.Tag {
	return tag.NewClosed("section", map[string]string{
		"class": "transclusion transclusion_failed transclusion_not-exists",
	}, []lines.Line{},
		tag.NewClosed("p", map[string]string{}, []lines.Line{lines.IndentableFrom(errParagraph)}))
}

// MapTransclusionErrorToTag returns an error tag that you should display to the user. If there is no error in the transclusion, bad things will happen, so verify with xcl.HasError beforehand.
func MapTransclusionErrorToTag(xcl blocks.Transclusion) tag.Tag {
	switch xcl.TransclusionError.Reason {
	case blocks.TransclusionErrorNotExists:
		return wrapInTransclusionError(fmt.Sprintf(`Cannot transclude hypha <a class="wikilink wikilink_new" href="/hypha/%[1]s">%[1]s</a> because it does not exist`, xcl.Target))

	case blocks.TransclusionErrorNoTarget:
		return wrapInTransclusionError("Transclusion target not specified")
	case blocks.TransclusionInTerminal:
		return wrapInTransclusionError("Transclusion is not supported in documents generated using Mycomarkup CLI")
	case blocks.TransclusionErrorOldSyntax:
		return wrapInTransclusionError("This transclusion is using the old syntax. Please update it to the new one")
	}
	log.Printf("MapTransclusionErrorToTag: unknown kind of transclusion error %d\n", xcl.TransclusionError.Reason)
	return tag.Tag{}
}
