package httpproxy

import (
	"path"
	"strings"
)

type HostMatcher struct {
	starValue    interface{}
	strictMap    map[string]interface{}
	prefixList   []string
	prefixMap    map[string]interface{}
	wildcardList []string
	wildcardMap  map[string]interface{}
}

func (hm *HostMatcher) add(host string, value interface{}) {
	switch {
	case host == "*":
		hm.starValue = value
	case !strings.Contains(host, "*"):
		hm.strictMap[host] = value
	case strings.HasPrefix(host, "*") && !strings.Contains(host[1:], "*"):
		hm.prefixList = append(hm.prefixList, host[1:])
		hm.prefixMap[host[1:]] = value
	default:
		hm.wildcardList = append(hm.wildcardList, host)
		hm.wildcardMap[host] = value
	}
}

func NewHostMatcher(hosts []string) *HostMatcher {
	values := make(map[string]interface{}, len(hosts))
	for _, host := range hosts {
		values[host] = struct{}{}
	}
	return NewHostMatcherWithValue(values)
}

func NewHostMatcherWithString(hosts map[string]string) *HostMatcher {
	values := make(map[string]interface{}, len(hosts))
	for host, value := range hosts {
		values[host] = value
	}
	return NewHostMatcherWithValue(values)
}

func NewHostMatcherWithValue(values map[string]interface{}) *HostMatcher {
	hm := &HostMatcher{
		starValue:    nil,
		strictMap:    make(map[string]interface{}),
		prefixList:   make([]string, 0),
		prefixMap:    make(map[string]interface{}),
		wildcardList: make([]string, 0),
		wildcardMap:  make(map[string]interface{}),
	}

	for host, value := range values {
		hm.add(host, value)
	}

	return hm
}

func (hm *HostMatcher) AddHost(host string) {
	hm.AddHostWithValue(host, struct{}{})
}

func (hm *HostMatcher) AddHostWithValue(host string, value interface{}) {
	hm.add(host, value)
}

func (hm *HostMatcher) Match(host string) bool {
	_, ok := hm.Lookup(host)
	return ok
}

func (hm *HostMatcher) Lookup(host string) (interface{}, bool) {
	if hm.starValue != nil {
		return hm.starValue, true
	}

	if value, ok := hm.strictMap[host]; ok {
		return value, true
	}

	for _, key := range hm.prefixList {
		if strings.HasSuffix(host, key) {
			return hm.prefixMap[key], true
		}
	}

	for _, key := range hm.wildcardList {
		if matched, _ := path.Match(key, host); matched {
			return hm.wildcardMap[key], true
		}
	}

	return nil, false
}
