package main

import "github.com/Reticent93/snips/pkg/models"

type templateData struct {
	Snip  *models.Snip   `json:"snip,omitempty"`
	Snips []*models.Snip `json:"snips,omitempty"`
}
