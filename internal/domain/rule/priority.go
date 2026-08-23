package rule

import "sort"

type ByPriority []Rule

func (r ByPriority) Len() int { return len(r) }
func (r ByPriority) Less(i, j int) bool {
	if r[i].EscalationMinutes == r[j].EscalationMinutes {
		return r[i].CreatedAt.Before(r[j].CreatedAt)
	}
	return r[i].EscalationMinutes < r[j].EscalationMinutes
}
func (r ByPriority) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func Sort(rules []Rule) []Rule {
	out := make([]Rule, len(rules))
	copy(out, rules)
	sort.Sort(ByPriority(out))
	return out
}
func FirstMatch(rules []Rule, source, name string) (Rule, bool) {
	for _, r := range Sort(rules) {
		if r.Enabled && r.Source == source && r.EventName == name {
			return r, true
		}
	}
	return Rule{}, false
}
