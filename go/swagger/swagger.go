package swagger

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.yml
var spec []byte

const page = `<!doctype html>
<html><head>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui.css">
</head><body>
<div id="swagger"></div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
<script>SwaggerUIBundle({ url: "/openapi.yaml", dom_id: "#swagger" })</script>
</body></html>`

func Routes() {
	http.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(spec)
	})
	http.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(page))
	})
}
