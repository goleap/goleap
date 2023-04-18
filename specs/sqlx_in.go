package specs

type SqlIn func(query string, args ...interface{}) (string, []interface{}, error)
