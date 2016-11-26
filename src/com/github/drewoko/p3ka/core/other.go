package core

import "strings"

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if strings.ToLower(a) == strings.ToLower(e) {
			return true
		}
	}
	return false
}

type Msg struct {
	Id int64
	Name string
	Text string
	Channel string
}

type Config struct {
	Database string
	Port string
	Static string
	BannedUsers []string
	ExcludedUsers []string
	Peka2TvHost string
	Peka2TvPort int
}