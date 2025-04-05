package config

type Config struct {
	AppPort       string
	TracePort     string
	AuthUser      string
	AuthPassword  string
	AuthJwtSecret string
	DbPath        string
	RateLimit     float64
	RateBurst     int
}
