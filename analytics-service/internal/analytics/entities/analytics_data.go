package entities

type AnalyticsData struct {
	Events       []*AnalyticsEvent
	TotalClicks  uint32
	UniqueClicks uint32
	Page         uint32
	PerPage      uint32
	TotalPages   uint32
}
