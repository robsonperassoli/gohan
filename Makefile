templ:
	@templ generate -watch -proxy=http://localhost:4100

tailwind:
	@tailwindcss -i views/css/app.css -o public/styles.css --watch
