package main

import (
	v1 "yuki-kabe/api/yuki-kabe/v1"
	"yuki-kabe/internal/data"
)

func main() {
	_ = data.ProviderSet
	_ = v1.Version
}
