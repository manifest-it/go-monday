package simpleitem

import (
	"strings"
	"time"

	"github.com/grokify/go-monday"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/stringsutil"
)

const (
	NoStatus = "NO STATUS"
	DONE     = "DONE"
	BLOCKED  = "BLOCKED"
	TBD      = "TBD"
	WIP      = "WIP"
)

type SimpleItem struct {
	Name                  string
	Status                string
	Date                  *time.Time
	LastChangedStatusDate *time.Time
	FieldsSimple          []SimpleCell
}

type SimpleCell struct {
	Title     string
	ChangedAt *time.Time
	Value     string
}

func (si *SimpleItem) TrimSpace() {
	si.Name = strings.TrimSpace(si.Name)
	si.Status = strings.TrimSpace(si.Status)
}

func (si *SimpleItem) String(inclStatus bool) string {
	var parts []string
	si.TrimSpace()
	if len(si.Name) > 0 {
		parts = append(parts, si.Name+":")
	}
	if si.Date == nil || timeutil.IsZero(*si.Date) {
		parts = append(parts, TBD)
	} else {
		parts = append(parts, si.Date.Format(timeutil.MonthDay))
	}
	if inclStatus {
		siStatus := stringsutil.TrimSpaceOrDefault(
			strings.ToUpper(si.Status), NoStatus)
		parts = append(parts, "("+siStatus+")")
	}
	if si.LastChangedStatusDate != nil && !timeutil.IsZero(*si.LastChangedStatusDate) {
		parts = append(parts, "(updated: "+si.LastChangedStatusDate.Format(timeutil.MonthDay)+")")
	}

	return strings.Join(parts, " ")
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
		si.LastChangedStatusDate = &changedAt
	}
	return si
}
