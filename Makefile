live/templ:
	templ generate --watch --proxy="http://localhost:4269" --open-browser=false -v


# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/air-verse/air@latest \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit "/usr/bin/true"

# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	bun x @tailwindcss/cli -i ./input.css -o ./assets/styles.css --minify --watch

# run esbuild to generate the index.js bundle in watch mode.
live/esbuild:
	bun x esbuild js/index.ts --bundle --outdir=assets/
	bun x esbuild js/index.ts --bundle --outdir=assets/ --watch

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	go run github.com/air-verse/air@latest \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "/usr/bin/true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

build/tailwind:
	bun x @tailwindcss/cli -i ./input.css -o ./assets/styles.css --minify

build/templ:
	templ generate -v

build:
	bun i
	make -j2 build/templ build/tailwind
	curl -SL -o assets/htmx.min.js https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js
	find assets/ -type f \( -name '*.css' -o -name '*.js' \) -exec gzip -9 -k --force {} +

# start all 5 watch processes in parallel.
live: 
	make build
	make -j5 live/templ live/server live/tailwind live/sync_assets

