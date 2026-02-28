package constants

const (
	ProgramIdentifier = "waypoint"
	AssetDir          = "./data"
	UIAddress         = "https://localhost:5173" // This is the UI dev server address. Used to redirect users to the UI after OIDC login.
)

var (
	Version        = "unknown"
	BuildTimestamp = "unknown"
	CommitHash     = "unknown"
)
