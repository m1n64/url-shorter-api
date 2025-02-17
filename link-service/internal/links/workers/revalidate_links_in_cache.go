package workers

import (
	"fmt"
	"link-service/internal/links/services"
	"time"
)

func RevalidateLinksInCache(linkService *services.LinkService) {
	ticker := time.NewTicker(time.Minute * 30)
	defer ticker.Stop()

	for {
		fmt.Println("Starting cache links globally...")

		_ = linkService.SaveAllInGlobalCache()
		fmt.Println("Links successfully globally cached!")

		<-ticker.C
	}
}
