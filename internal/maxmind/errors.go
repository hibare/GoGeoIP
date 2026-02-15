package maxmind

import "errors"

var (
	// ErrCountryDBNotLoaded is returned when the country database is not loaded.
	ErrCountryDBNotLoaded = errors.New("country database not loaded")

	// ErrCityDBNotLoaded is returned when the city database is not loaded.
	ErrCityDBNotLoaded = errors.New("city database not loaded")

	// ErrASNDBNotLoaded is returned when the ASN database is not loaded.
	ErrASNDBNotLoaded = errors.New("ASN database not loaded")

	// ErrLicenseKeyRequired is returned when the MaxMind license key is missing.
	ErrLicenseKeyRequired = errors.New("MAXMIND_LICENSE_KEY is required")

	// ErrChecksumMismatch is returned when the downloaded file checksum does not match.
	ErrChecksumMismatch = errors.New("checksum mismatch")

	// ErrEmptySHA256File is returned when the SHA256 file is empty.
	ErrEmptySHA256File = errors.New("empty SHA256 file")

	// ErrInvalidSHA256File is returned when the SHA256 file format is invalid.
	ErrInvalidSHA256File = errors.New("invalid SHA256 file format")

	// ErrDBDownloadFailed is returned when downloading databases fails and no local copies exist.
	ErrDBDownloadFailed = errors.New("failed to download databases and no existing files found")

	// ErrDBOpenFailed is returned when opening a database file fails.
	ErrDBOpenFailed = errors.New("failed to open database")
)
