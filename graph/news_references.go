package graph

import (
	"encoding/json"

	"github.com/weeb-vip/anime-api/graph/model"
)

// decodeNewsReferences unpacks the anime_news.reference_links JSON column into the
// GraphQL type. Malformed JSON yields no references rather than failing the whole news
// query — the story is worth more than its attachments.
//
// Lives in its own file because gqlgen relocates any non-resolver function it finds in
// types.resolvers.go into a "was going to be deleted" block on every regeneration.
func decodeNewsReferences(raw *string) []*model.NewsReference {
	if raw == nil || *raw == "" {
		return nil
	}
	var decoded []struct {
		Kind  string `json:"kind"`
		Title string `json:"title"`
		URL   string `json:"url"`
	}
	if err := json.Unmarshal([]byte(*raw), &decoded); err != nil {
		return nil
	}
	var out []*model.NewsReference
	for _, d := range decoded {
		if d.URL == "" {
			continue
		}
		out = append(out, &model.NewsReference{Kind: d.Kind, Title: d.Title, URL: d.URL})
	}
	return out
}
