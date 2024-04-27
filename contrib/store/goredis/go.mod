module github.com/fzdwx/authz/contrib/store/goredis

go 1.18

require (
	github.com/fzdwx/authz v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.5.1
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.39.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/fzdwx/authz => ../../../
