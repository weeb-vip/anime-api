package anime

import (
	"strings"
)

// FieldSelection represents which fields should be selected in the query
type FieldSelection struct {
	Fields map[string]bool
}

// NewFieldSelection creates a new FieldSelection from a list of field names
func NewFieldSelection(fields []string) *FieldSelection {
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}
	return &FieldSelection{Fields: fieldMap}
}

// Has checks if a field is selected
func (fs *FieldSelection) Has(field string) bool {
	if fs == nil || fs.Fields == nil {
		return true // If no field selection, select all
	}
	return fs.Fields[field]
}

// BuildSelectClause builds a SQL SELECT clause based on the selected fields
func (fs *FieldSelection) BuildSelectClause(tableName string) string {
	if fs == nil || fs.Fields == nil || len(fs.Fields) == 0 {
		return tableName + ".*" // Select all if no specific fields
	}

	// Map GraphQL field names to database column names
	fieldMapping := map[string]string{
		"id":            "id",
		"anidbid":       "anidb_id",
		"thetvdbid":     "the_tvdb_id",
		"titleEn":       "title_en",
		"titleJp":       "title_jp",
		"titleRomaji":   "title_romaji",
		"titleKanji":    "title_kanji",
		"titleSynonyms": "title_synonyms",
		"description":   "synopsis",
		"imageUrl":      "image_url",
		"tags":          "genres",
		"studios":       "studios",
		"animeStatus":   "status",
		"episodeCount":  "episodes",
		"duration":      "duration",
		"rating":        "rating",
		"startDate":     "start_date",
		"endDate":       "end_date",
		"broadcast":     "broadcast",
		"source":        "source",
		"licensors":     "licensors",
		"ranking":       "ranking",
		"createdAt":     "created_at",
		"updatedAt":     "updated_at",
	}

	var selectedColumns []string

	// Always include ID for entity consistency
	selectedColumns = append(selectedColumns, tableName+".id")

	for field := range fs.Fields {
		if column, exists := fieldMapping[field]; exists && column != "id" {
			selectedColumns = append(selectedColumns, tableName+"."+column)
		}
	}

	// Always include created_at and updated_at for entity consistency
	if !fs.Has("createdAt") {
		selectedColumns = append(selectedColumns, tableName+".created_at")
	}
	if !fs.Has("updatedAt") {
		selectedColumns = append(selectedColumns, tableName+".updated_at")
	}

	if len(selectedColumns) == 0 {
		return tableName + ".*"
	}

	return strings.Join(selectedColumns, ", ")
}