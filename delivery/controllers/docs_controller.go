package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DocsController serves Swagger UI
type DocsController struct{}

func NewDocsController() *DocsController { return &DocsController{} }

// SwaggerUI serves a minimal Swagger UI page backed by /docs/openapi.yaml
func (d *DocsController) SwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>RemedyMate API Docs</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    html, body { margin: 0; padding: 0; height: 100%; }
    #swagger-ui { height: 100%; }
  </style>
  <link rel="icon" href="data:,">
  <!-- Intentionally using CDN for a small footprint -->
  <!-- To self-host, replace with embedded swagger-ui assets -->
  <script defer src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script defer src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.addEventListener('DOMContentLoaded', function () {
      window.ui = SwaggerUIBundle({
        url: '/api/v1/docs/openapi.yaml',
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
        layout: 'BaseLayout',
        deepLinking: true,
        syntaxHighlight: { activate: true },
      });
    });
  </script>
  </head>
<body>
  <div id="swagger-ui"></div>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
