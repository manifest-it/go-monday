package simpleitem

import (
	"sort"
	"time"

	"github.com/grokify/go-monday"
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
func (s SimpleItems) Sort() { sort.Sort(s) }

type SimpleItem struct {
	Name                    string
	Status                  string
	Date                    *time.Time
	LastChangedAtStatusDate *time.Time
	FieldsSimple            []SimpleCell
}

type SimpleCell struct {
	Title     string
	ChangedAt *time.Time
	Value     string
}

func ColumnValueToSimple(colvals []monday.ColumnValue) []SimpleCell {
	svals := []SimpleCell{}
	for _, cv := range colvals {
		sv := SimpleCell{
			Title: cv.Title}
		if cv.Text == nil {
			sv.Value = ""
		} else {
			sv.Value = *cv.Text
		}
		if cv.Value != nil && len(*cv.Value) > 0 {
			cvv, err := monday.ParseColumnValueValue([]byte(*cv.Value))
			if err == nil {
				sv.ChangedAt = cvv.ChangedAt
			}
		}

		svals = append(svals, sv)
	}
	return svals
}

func ItemToSimple(item monday.Item) SimpleItem {
	si := SimpleItem{
		Name:         item.Name,
		FieldsSimple: ColumnValueToSimple(item.ColumnValues)}
	date, err := item.Date()
	if err == nil {
		si.Date = &date
	}
	statusCv, err := item.GetColumnValue("Status", true)
	if err == nil && statusCv.Text != nil {
		si.Status = *statusCv.Text
	}
	changedAt, err := item.LastChangedAtDateStatus()
	if err == nil {
		si.LastChangedAtStatusDate = &changedAt
	}
	return si
}

func BoardSimpleItems(b monday.Board) (SimpleItems, error) {
	var simps SimpleItems
	for _, item := range b.Items {
		simps = append(simps, ItemToSimple(item))
	}
	return simps, nil
}
