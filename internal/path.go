package internal

import (
	"os"
)

func GetSiteUrl() string {
	return os.Getenv("SITE_URL")
}
