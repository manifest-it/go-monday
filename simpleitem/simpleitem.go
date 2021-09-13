package simpleitem

import (
	"strings"
	"time"

	"github.com/grokify/go-monday"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/grokify/simplego/time/timeutil"
	"github.com/grokify/simplego/type/stringsutil"
)

const (
	NoStatus = "NO STATUS"
	DONE     = "DONE"
	BLOCKED  = "BLOCKED"
	TBD      = "TBD"
	WIP      = "WIP"

	BulletTypeNumeric = "numeric"
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
	LinkURL   string
	LinkText  string
}

func (si *SimpleItem) TrimSpace() {
	si.Name = strings.TrimSpace(si.Name)
	si.Status = strings.TrimSpace(si.Status)
}

func (si *SimpleItem) LinkURL() string {
	for _, sc := range si.FieldsSimple {
		sc.LinkURL = strings.TrimSpace(sc.LinkURL)
		if len(sc.LinkURL) > 0 {
			return sc.LinkURL
		}
	}
	return ""
}

func (si *SimpleItem) String(linkify, inclStatus bool) string {
	var parts []string
	si.TrimSpace()
	if len(si.Name) > 0 {
		parts = append(parts, si.Name+":")
	}
	if si.Date == nil || timeutil.IsZero(*si.Date) {
		parts = append(parts, TBD)
	} else if si.Date.UTC().Year() == time.Now().UTC().Year() {
		parts = append(parts, si.Date.Format(timeutil.MonthDay))
	} else {
		parts = append(parts, si.Date.Format("_1/_2/06"))
	}
	if inclStatus {
		siStatus := stringsutil.TrimSpaceOrDefault(
			strings.ToUpper(si.Status), NoStatus)
		parts = append(parts, "("+siStatus+")")
	}
	if si.LastChangedStatusDate != nil && !timeutil.IsZero(*si.LastChangedStatusDate) {
		parts = append(parts, "(updated: "+si.LastChangedStatusDate.Format(timeutil.MonthDay)+")")
	}

	text := strings.Join(parts, " ")
	if linkify {
		linkURL := strings.TrimSpace(si.LinkURL())
		if len(linkURL) > 0 {
			return "[" + text + "](" + linkURL + ")"
		}
	}
	return text
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
				if cv.Id == monday.ColumnValueIdLink {
					cvv.URL = strings.TrimSpace(cvv.URL)
					cvv.Text = strings.TrimSpace(cvv.Text)
					if urlutil.IsHttp(cvv.URL, true, true) {
						sv.LinkURL = cvv.URL
						sv.LinkText = cvv.Text
					}
				}
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
	statusCv, err := item.GetColumnValue(monday.ColumnValueTitleStatus, true)
	if err == nil && statusCv.Text != nil {
		si.Status = *statusCv.Text
	}
	changedAt, err := item.LastChangedAtDateStatus()
	if err == nil {
		si.LastChangedStatusDate = &changedAt
	}
	return si
}
