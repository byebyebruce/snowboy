module github.com/Kitt-AI/snowboy

go 1.17

require (
	github.com/brentnd/go-snowboy v0.0.0-20190301212623-e19133c572af
	github.com/gordonklaus/portaudio v0.0.0-20200911161147-bb74aa485641
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/brentnd/go-snowboy v0.0.0-20190301212623-e19133c572af => ./go-snowboy
