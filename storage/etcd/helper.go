package etcd

import (
	"strings"
)

const (
	trace = "trace"
	tags  = "tags"
	ltags = "tags/"
	sep   = "/"
)

/*

Trace key examples:
trace/123sdfgdfhgfg/1 -> span with id 1
trace/123sdfgdfhgfg/1/1 -> child span with id 1 of span with id 1
trace/123sdfgdfhgfg/1/2
trace/123sdfgdfhgfg/2

Trace tags query key example:
trace/tags/123sdfgdfhgfg/1 -> span with id 1
trace/tags/123sdfgdfhgfg/1/1 -> child span with id 1 of span with id 1
trace/tags/123sdfgdfhgfg/1/2
trace/tags/123sdfgdfhgfg/2

*/

func traceKey(traceId string) string {
	s := []string{trace, traceId}
	return strings.Join(s, "/")
}

func formTagKey(parts []string) string {
	s := []string{trace, tags}
	return strings.Join(append(s, parts...), "/")
}

func formKey(parts []string) string {
	s := []string{trace}
	return strings.Join(append(s, parts...), "/")
}

func lookupKey() string {
	s := []string{trace, tags}
	return strings.Join(s, "/")
}

func transformKey(lookupKey string) string {
	return strings.Replace(lookupKey, ltags, "", -1)
}

func extractTraceKey(key string) string {
	keyTransformed := transformKey(key)
	return strings.Split(keyTransformed, sep)[1]
}
