package adb

type JQL struct {
	ops map[string]func()
}

func (jql *JQL) Query(sql string, json []byte) {
	//ret := regexp.MustCompile(`/^(select)\s+([a-z0-9_\,\.\s\*]+)\s+from\s+([a-z0-9_\.]+)(?: where\s+\((.+)\))?\s*(?:order\sby\s+([a-z0-9_\,]+))?\s*(asc|desc|ascnum|descnum)?\s*(?:limit\s+([0-9_\,]+))?/i`)
	// NOTE: finish this...
	//
	return jql.Parse(json, ops)
}

func (jql *JQL) Parse(json []byte) {

}
