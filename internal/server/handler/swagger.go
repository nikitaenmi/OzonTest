package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) ServeSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	swaggerPath := filepath.Join("gen", "payment.swagger.json")
	
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		http.Error(w, "Swagger JSON not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *SwaggerHandler) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>OzonTest Payment API</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: '/swagger.json',
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.presets.debug
            ],
            layout: "BaseLayout"
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}