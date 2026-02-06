// NOTE: This module path is temporarily set to github.com/dxcSithLord for development.
// Revert to github.com/sqrldev/server-go-ssp-gormauthstore before pushing to upstream.
module github.com/dxcSithLord/server-go-ssp-gormauthstore

go 1.24.0

toolchain go1.24.7

require (
	github.com/dxcSithLord/server-go-ssp v0.0.0-20260202110616-66529f78b7f1
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.33 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/yeqown/go-qrcode/v2 v2.2.5 // indirect
	github.com/yeqown/go-qrcode/writer/standard v1.3.0 // indirect
	github.com/yeqown/reedsolomon v1.0.0 // indirect
	golang.org/x/image v0.35.0 // indirect
	golang.org/x/text v0.33.0 // indirect
)

// NOTE: Replace directives for local development.
// Uncomment and adjust paths as needed for local builds.
// Remove or comment out before pushing to upstream.
// replace github.com/dxcSithLord/server-go-ssp => ../server-go-ssp
