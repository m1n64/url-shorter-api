package workers

import (
	"fmt"
	"go.uber.org/zap"
	"time"
	"user-service/internal/users/repositories"
)

func StartRemoveExpiredTokensWorker(repo repositories.TokenRepository, logger *zap.Logger, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		logger.Info(fmt.Sprintf("Removing expired tokens from last %s...", interval.String()))

		err := repo.DeleteExpiredTokens()
		if err != nil {
			logger.Error("Error removing expired tokens", zap.Error(err))
		} else {
			logger.Info("Expired tokens removed successfully.")
		}
	}
}
