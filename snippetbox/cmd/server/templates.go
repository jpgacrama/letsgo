package server

import "snippetbox/pkg/models"

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
