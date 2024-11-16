package netconfigparser

import "github.com/opennetworktools/NetConfigParser/internal/parser"

func GetParser(osType, path string) parser.ConfigParser {
	switch osType {
	case "IOSXE":
		return parser.NewIOSXEParser(path)
	default:
		return nil
	}
}
