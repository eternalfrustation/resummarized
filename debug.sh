#!/bin/bash

templ generate --watch --open-browser=false &
pnpm tailwindcss -i ./input.css -o ./assets/styles.css --minify --watch &
air
