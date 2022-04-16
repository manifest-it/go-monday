package simpleitem

import (
	"sort"
	"strconv"
	"strings"

	monday "github.com/grokify/go-monday"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	SortAsc  = "asc"
	SortDec  = "dec"
	SortDesc = "desc"
	SortDsc  = "dsc"
)

type SimpleItems []SimpleItem

func (s SimpleItems) Len() int { return len(s) }
func (s SimpleItems) Less(i, j int) bool {
	if s[i].Date == nil && s[j].Date == nil {
		return false
	} else if s[i].Date == nil {
		return false
	} else if s[j].Date == nil {
		return true
	}
	return s[i].Date.Before(*s[j].Date)
}
func (s SimpleItems) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort is a convenience method.
func (s SimpleItems) Sort(reverse bool) {
	if reverse {
		sort.Sort(sort.Reverse(s))
	} else {
		sort.Sort(s)
	}
}

func (s SimpleItems) ByStatus() map[string]SimpleItems {
	byStatus := map[string]SimpleItems{}
	for _, si := range s {
		statusUpper := stringsutil.TrimSpaceOrDefault(
			strings.ToUpper(si.Status), NoStatus)
		if _, ok := byStatus[statusUpper]; !ok {
			byStatus[statusUpper] = SimpleItems{}
		}
		byStatus[statusUpper] = append(byStatus[statusUpper], si)
	}
	return byStatus
}

func (s SimpleItems) Strings(bulletType, delimit, sortType string, linkify, inclStatus bool) []string {
	bulletType = strings.ToLower(strings.TrimSpace(bulletType))
	lines := []string{}
	sortType = strings.ToLower(strings.TrimSpace(sortType))
	if sortType == SortAsc {
		s.Sort(false)
	} else if sortType == SortDec || sortType == SortDesc || sortType == SortDsc {
		s.Sort(true)
	}
	for i, item := range s {
		str := item.String(linkify, inclStatus)
		if bulletType == BulletTypeNumeric {
			str = strconv.Itoa(i+1) + delimit + str
		} else if len(bulletType) > 0 {
			str = bulletType + delimit + str
		}
		lines = append(lines, str)
	}
	return lines
}

func (s SimpleItems) Statuses(convert bool, skip []string) []string {
	statuses := []string{}
	skipMap := map[string]int{}
	for _, skipStatus := range skip {
		skipMap[skipStatus] = 1
	}
	for _, item := range s {
		status := item.Status
		if convert {
			status = strings.ToUpper(strings.TrimSpace(status))
			if len(status) == 0 {
				status = NoStatus
			}
		}
		if _, ok := skipMap[status]; !ok {
			statuses = append(statuses, status)
		}
	}
	sort.Strings(statuses)
	return stringsutil.SliceCondenseSpace(statuses, true, true)
}

func (s SimpleItems) StringsByStatus(bulletType, delimit, sortType string, linkify, inclStatus bool) []string {
	lines := []string{}
	byStatus := s.ByStatus()
	stdStatuses := []string{DONE, WIP, BLOCKED, NoStatus}
	othStatuses := s.Statuses(true, stdStatuses)
	allStatuses := append(stdStatuses, othStatuses...)
	for _, statusName := range allStatuses {
		if items, ok := byStatus[statusName]; ok {
			lines = append(lines, statusName)
			itemsLines := items.Strings(bulletType, delimit, sortType, linkify, inclStatus)
			lines = append(lines, itemsLines...)
		}
	}
	return lines
}

func BoardSimpleItems(b monday.Board) (SimpleItems, error) {
	var simps SimpleItems
	for _, item := range b.Items {
		simps = append(simps, ItemToSimple(item))
	}
	return simps, nil
}
