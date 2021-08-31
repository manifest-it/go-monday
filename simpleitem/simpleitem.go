package simpleitem

import (
	"time"

	"github.com/grokify/go-monday"
)

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

func BoardSimpleItems(b monday.Board) ([]SimpleItem, error) {
	simps := []SimpleItem{}
	for _, item := range b.Items {
		simps = append(simps, ItemToSimple(item))
	}
	return simps, nil
}
