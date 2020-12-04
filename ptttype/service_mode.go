package ptttype

type ServiceMode int

const (
	_          ServiceMode = iota
	DEV                    // Development Mode
	PRODUCTION             // Production Mode
	DEBUG                  // Debug Mode (in production)
)

func (s ServiceMode) String() string {
	switch s {
	case DEV:
		return "DEV"
	case PRODUCTION:
		return "PRODUCTION"
	case DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

func stringToServiceMode(str string) ServiceMode {
	switch str {
	case "DEV":
		return DEV
	case "PRODUCTION":
		return PRODUCTION
	case "DEBUG":
		return DEBUG
	default:
		return DEV
	}
}
