package date

import (
	"fmt"
	"time"
)

// commonTimezoneAbbreviations maps common (sometimes ambiguous) abbreviations to
// a canonical IANA timezone name. Where an abbreviation is genuinely ambiguous,
// the most-used zone is preferred.
var commonTimezoneAbbreviations = map[string]string{
	"ACDT":  "Australia/Adelaide",
	"ACST":  "Australia/Darwin",
	"ADT":   "America/Halifax",
	"AEDT":  "Australia/Sydney",
	"AEST":  "Australia/Brisbane",
	"AFT":   "Asia/Kabul",
	"AKDT":  "America/Anchorage",
	"AKST":  "America/Anchorage",
	"ALMT":  "Asia/Almaty",
	"AMST":  "America/Manaus",
	"AMT":   "America/Manaus",
	"ANAST": "Asia/Anadyr",
	"ANAT":  "Asia/Anadyr",
	"AQTT":  "Asia/Aqtau",
	"ART":   "America/Argentina/Buenos_Aires",
	"AST":   "America/Halifax",
	"AWDT":  "Australia/Perth",
	"AWST":  "Australia/Perth",
	"AZOST": "Atlantic/Azores",
	"AZOT":  "Atlantic/Azores",
	"AZST":  "Asia/Baku",
	"AZT":   "Asia/Baku",
	"BNT":   "Asia/Brunei",
	"BOT":   "America/La_Paz",
	"BRST":  "America/Sao_Paulo",
	"BRT":   "America/Sao_Paulo",
	"BST":   "Europe/London",
	"BTT":   "Asia/Thimphu",
	"CAST":  "Antarctica/Casey",
	"CAT":   "Africa/Harare",
	"CCT":   "Indian/Cocos",
	"CDT":   "America/Chicago",
	"CEST":  "Europe/Paris",
	"CET":   "Europe/Paris",
	"CHADT": "Pacific/Chatham",
	"CHAST": "Pacific/Chatham",
	"CKT":   "Pacific/Rarotonga",
	"CLST":  "America/Santiago",
	"CLT":   "America/Santiago",
	"COT":   "America/Bogota",
	"CST":   "America/Chicago",
	"CVT":   "Atlantic/Cape_Verde",
	"CXT":   "Indian/Christmas",
	"ChST":  "Pacific/Guam",
	"DAVT":  "Antarctica/Davis",
	"EASST": "Pacific/Easter",
	"EAST":  "Pacific/Easter",
	"EAT":   "Africa/Nairobi",
	"ECT":   "America/Guayaquil",
	"EDT":   "America/New_York",
	"EEST":  "Europe/Helsinki",
	"EET":   "Europe/Helsinki",
	"EGST":  "America/Scoresbysund",
	"EGT":   "America/Scoresbysund",
	"EST":   "America/New_York",
	"ET":    "America/New_York",
	"FJST":  "Pacific/Fiji",
	"FJT":   "Pacific/Fiji",
	"FKST":  "Atlantic/Stanley",
	"FKT":   "Atlantic/Stanley",
	"FNT":   "America/Noronha",
	"GALT":  "Pacific/Galapagos",
	"GAMT":  "Pacific/Gambier",
	"GET":   "Asia/Tbilisi",
	"GFT":   "America/Cayenne",
	"GILT":  "Pacific/Tarawa",
	"GMT":   "Etc/GMT",
	"GST":   "Asia/Dubai",
	"GYT":   "America/Guyana",
	"HAA":   "America/Halifax",
	"HAC":   "America/Winnipeg",
	"HADT":  "America/Adak",
	"HAE":   "America/New_York",
	"HAP":   "America/Vancouver",
	"HAR":   "America/Denver",
	"HAST":  "America/Adak",
	"HAT":   "America/St_Johns",
	"HAY":   "America/Anchorage",
	"HKT":   "Asia/Hong_Kong",
	"HLV":   "America/Caracas",
	"HNA":   "America/Halifax",
	"HNC":   "America/Winnipeg",
	"HNE":   "America/New_York",
	"HNP":   "America/Vancouver",
	"HNR":   "America/Denver",
	"HNT":   "America/St_Johns",
	"HNY":   "America/Anchorage",
	"HOVT":  "Asia/Hovd",
	"HST":   "Pacific/Honolulu",
	"ICT":   "Asia/Bangkok",
	"IDT":   "Asia/Jerusalem",
	"IOT":   "Indian/Chagos",
	"IRDT":  "Asia/Tehran",
	"IRKST": "Asia/Irkutsk",
	"IRKT":  "Asia/Irkutsk",
	"IRST":  "Asia/Tehran",
	"IST":   "Asia/Kolkata",
	"JST":   "Asia/Tokyo",
	"KGT":   "Asia/Bishkek",
	"KRAST": "Asia/Krasnoyarsk",
	"KRAT":  "Asia/Krasnoyarsk",
	"KST":   "Asia/Seoul",
	"KUYT":  "Europe/Samara",
	"LHDT":  "Australia/Lord_Howe",
	"LHST":  "Australia/Lord_Howe",
	"LINT":  "Pacific/Kiritimati",
	"MAGST": "Asia/Magadan",
	"MAGT":  "Asia/Magadan",
	"MART":  "Pacific/Marquesas",
	"MAWT":  "Antarctica/Mawson",
	"MDT":   "America/Denver",
	"MEST":  "Europe/Paris",
	"MET":   "Europe/Paris",
	"MHT":   "Pacific/Majuro",
	"MMT":   "Asia/Rangoon",
	"MSD":   "Europe/Moscow",
	"MSK":   "Europe/Moscow",
	"MST":   "America/Denver",
	"MUT":   "Indian/Mauritius",
	"MVT":   "Indian/Maldives",
	"MYT":   "Asia/Kuala_Lumpur",
	"NCT":   "Pacific/Noumea",
	"NDT":   "America/St_Johns",
	"NFT":   "Pacific/Norfolk",
	"NOVST": "Asia/Novosibirsk",
	"NOVT":  "Asia/Novosibirsk",
	"NPT":   "Asia/Kathmandu",
	"NST":   "America/St_Johns",
	"NUT":   "Pacific/Niue",
	"NZDT":  "Pacific/Auckland",
	"NZST":  "Pacific/Auckland",
	"OMSST": "Asia/Omsk",
	"OMST":  "Asia/Omsk",
	"ORAT":  "Asia/Oral",
	"PDT":   "America/Los_Angeles",
	"PET":   "America/Lima",
	"PETST": "Asia/Kamchatka",
	"PETT":  "Asia/Kamchatka",
	"PGT":   "Pacific/Port_Moresby",
	"PHOT":  "Pacific/Enderbury",
	"PHT":   "Asia/Manila",
	"PKT":   "Asia/Karachi",
	"PMDT":  "America/Miquelon",
	"PMST":  "America/Miquelon",
	"PONT":  "Pacific/Pohnpei",
	"PST":   "America/Los_Angeles",
	"PT":    "America/Los_Angeles",
	"PWT":   "Pacific/Palau",
	"PYST":  "America/Asuncion",
	"PYT":   "America/Asuncion",
	"QYZT":  "Asia/Qyzylorda",
	"RET":   "Indian/Reunion",
	"ROTT":  "Antarctica/Rothera",
	"SAKST": "Asia/Sakhalin",
	"SAKT":  "Asia/Sakhalin",
	"SAMT":  "Europe/Samara",
	"SAST":  "Africa/Johannesburg",
	"SBT":   "Pacific/Guadalcanal",
	"SCT":   "Indian/Mahe",
	"SGT":   "Asia/Singapore",
	"SRT":   "America/Paramaribo",
	"SST":   "Pacific/Pago_Pago",
	"SYOT":  "Antarctica/Syowa",
	"TAHT":  "Pacific/Tahiti",
	"TFT":   "Indian/Kerguelen",
	"TJT":   "Asia/Dushanbe",
	"TKT":   "Pacific/Fakaofo",
	"TLT":   "Asia/Dili",
	"TMT":   "Asia/Ashgabat",
	"TOT":   "Pacific/Tongatapu",
	"TVT":   "Pacific/Funafuti",
	"ULAT":  "Asia/Ulaanbaatar",
	"UTC":   "UTC",
	"UYST":  "America/Montevideo",
	"UYT":   "America/Montevideo",
	"UZT":   "Asia/Tashkent",
	"VET":   "America/Caracas",
	"VLAST": "Asia/Vladivostok",
	"VLAT":  "Asia/Vladivostok",
	"VOST":  "Antarctica/Vostok",
	"VUT":   "Pacific/Efate",
	"WAST":  "Africa/Windhoek",
	"WAT":   "Africa/Lagos",
	"WEST":  "Europe/Lisbon",
	"WET":   "Europe/Lisbon",
	"WFT":   "Pacific/Wallis",
	"WGST":  "America/Godthab",
	"WGT":   "America/Godthab",
	"WIB":   "Asia/Jakarta",
	"WIT":   "Asia/Jayapura",
	"WITA":  "Asia/Makassar",
	"WST":   "Pacific/Apia",
	"WT":    "Africa/Monrovia",
	"YAKST": "Asia/Yakutsk",
	"YAKT":  "Asia/Yakutsk",
	"YEKST": "Asia/Yekaterinburg",
	"YEKT":  "Asia/Yekaterinburg",
}

// IsValidTimezone reports whether zone is a recognized IANA timezone name.
func IsValidTimezone(zone string) bool {
	_, err := time.LoadLocation(zone)
	return err == nil
}

// LoadLocation is a convenience wrapper around time.LoadLocation that returns
// a descriptive error.
func LoadLocation(zone string) (*time.Location, error) {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return nil, fmt.Errorf("unknown timezone %q: %w", zone, err)
	}
	return loc, nil
}

// ConvertTZ converts t from its current timezone to the named IANA timezone.
func ConvertTZ(t time.Time, to string) (time.Time, error) {
	loc, err := LoadLocation(to)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// TimezoneOffset returns the UTC offset in seconds for the named IANA timezone
// at the moment represented by t.
func TimezoneOffset(zone string, t time.Time) (int, error) {
	loc, err := LoadLocation(zone)
	if err != nil {
		return 0, err
	}
	_, offset := t.In(loc).Zone()
	return offset, nil
}

// ParseInTimezone parses the date string s using dateparse.ParseAny and then
// reinterprets the resulting time as being in the named IANA timezone.
// Use this when the date string carries no timezone information.
func ParseInTimezone(s, zone string) (time.Time, error) {
	loc, err := LoadLocation(zone)
	if err != nil {
		return time.Time{}, err
	}
	return ParseIn(s, loc)
}

// LocalToUTC converts t to UTC.
func LocalToUTC(t time.Time) time.Time {
	return t.UTC()
}

// UTCToLocal converts a UTC time to the system local timezone.
func UTCToLocal(t time.Time) time.Time {
	return t.In(time.Local)
}

// GuessTimezone maps a common timezone abbreviation (e.g. "IST", "PST") to
// a canonical IANA timezone name. It returns the name and true on success, or
// ("", false) when the abbreviation is not recognized.
func GuessTimezone(abbr string) (string, bool) {
	name, ok := commonTimezoneAbbreviations[abbr]
	return name, ok
}
