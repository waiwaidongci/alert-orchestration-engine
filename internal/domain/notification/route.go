package notification

import "strings"

type Route struct {
	Channel string
	Target  string
	Enabled bool
	Labels  map[string]string
}

func (r Route) Matches(labels map[string]string) bool {
	if !r.Enabled {
		return false
	}
	for k, v := range r.Labels {
		if labels[k] != v {
			return false
		}
	}
	return true
}
func (r Route) Normalized() Route {
	r.Channel = strings.ToLower(strings.TrimSpace(r.Channel))
	r.Target = strings.TrimSpace(r.Target)
	if r.Labels == nil {
		r.Labels = map[string]string{}
	}
	return r
}
func SelectRoutes(routes []Route, labels map[string]string) []Route {
	out := routes[:0]
	for _, r := range routes {
		r = r.Normalized()
		if r.Matches(labels) {
			out = append(out, r)
		}
	}
	return out
}
