package navigator

import (
	"strings"
)

type selectors []selector

func (ss selectors) canMergeType(selectorType selectorType) bool {
	if len(ss) == 0 {
		return false
	}
	last := ss[len(ss)-1]
	bothCSS := selectorType == cssType && last.Type == cssType
	return bothCSS && !last.Indexed && !last.Single
}

// Append clones selectors and append (or merge) new selector.
func (ss selectors) Append(selectorType selectorType, value string) selectors {
	sl := selector{
		Type:  selectorType,
		Value: value,
	}
	if !ss.canMergeType(selectorType) {
		return ss.clonePlusOne(sl)
	}
	idx := len(ss) - 1
	sl.Value = ss[idx].Value + " " + sl.Value
	return ss[:idx].clonePlusOne(sl)
}

func (ss selectors) Single() selectors {
	idx := len(ss) - 1
	if idx < 0 {
		return nil
	}
	sl := ss[idx]
	sl.Single = true
	sl.Indexed = false
	return ss[:idx].clonePlusOne(sl)
}

func (ss selectors) At(index int) selectors {
	idx := len(ss) - 1
	if idx < 0 {
		return nil
	}
	sl := ss[idx]
	sl.Single = false
	sl.Indexed = true
	sl.Index = index
	return ss[:idx].clonePlusOne(sl)
}

func (ss selectors) String() string {
	var tags []string
	for _, selector := range ss {
		tags = append(tags, selector.String())
	}
	return strings.Join(tags, " | ")
}

func (ss selectors) clonePlusOne(one selector) selectors {
	ret := make(selectors, len(ss), len(ss)+1)
	copy(ret, ss)
	ret = append(ret, one)
	return ret
}
