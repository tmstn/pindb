package pindb

import "fmt"

type Warnings []Warning

func (w *Warnings) json() WarningsJSON {
	j := WarningsJSON{}
	for _, v := range *w {
		j = append(j, v.JSON())
	}
	return j
}

func newWarnings() Warnings {
	return Warnings{}
}

type WarningCategory string

func (w WarningCategory) String() string {
	return string(w)
}

func (w WarningCategory) MismatchRecord() bool {
	return w == MismatchRecordWarning
}

func (w WarningCategory) NoUUID() bool {
	return w == NoUUIDWarning
}

func (w WarningCategory) MismatchUUID() bool {
	return w == MismatchUUIDWarning
}

func (w WarningCategory) MultiplePinDBGroupTag() bool {
	return w == MultiplePinDBGroupTagWarning
}

func (w WarningCategory) UnrelatedPinDBGroupTag() bool {
	return w == UnrelatedPinDBGroupTagWarning
}

func (w WarningCategory) UnrelatedPinDBStoreTag() bool {
	return w == UnrelatedPinDBStoreTagWarning
}

func (w WarningCategory) UnrelatedPinDBBucketTag() bool {
	return w == UnrelatedPinDBBucketTagWarning
}

const (
	MismatchRecordWarning          WarningCategory = "mismatch_record"
	NoUUIDWarning                  WarningCategory = "no_uuid"
	MismatchUUIDWarning            WarningCategory = "mismatch_uuid"
	MultiplePinDBGroupTagWarning   WarningCategory = "multiple_pindb_group_tag"
	UnrelatedPinDBGroupTagWarning  WarningCategory = "unrelated_pindb_group_tag"
	UnrelatedPinDBStoreTagWarning  WarningCategory = "unrelated_pindb_store_tag"
	UnrelatedPinDBBucketTagWarning WarningCategory = "unrelated_pindb_bucket_tag"
)

type Warning struct {
	category WarningCategory
	tag      Tag
}

func (w *Warning) String() string {
	switch true {
	case w.category.MismatchRecord():
		return "the link description record does not match data"
	case w.category.MismatchUUID():
		return "the link uuid does not match the data"
	case w.category.MultiplePinDBGroupTag():
		return fmt.Sprintf("the link has multiple group tags (%s)", w.tag.String())
	case w.category.UnrelatedPinDBGroupTag():
		return fmt.Sprintf("the link has group tags outside of the current bucket (%s)", w.tag.String())
	case w.category.UnrelatedPinDBBucketTag():
		return fmt.Sprintf("the link has a bucket tag outside of the current bucket (%s)", w.tag.String())
	case w.category.UnrelatedPinDBStoreTag():
		return fmt.Sprintf("the link has a store tag outside of the current store (%s)", w.tag.String())
	default:
		return "unknown warning"
	}
}

func (w *Warning) JSON() WarningJSON {
	j := WarningJSON{}
	j.Category = w.category.String()
	j.Tag = w.tag.String()
	return j
}

func NewWarning(category WarningCategory, tag *Tag) Warning {
	w := Warning{
		category: category,
	}

	if tag != nil {
		w.tag = *tag
	}

	return w
}

type WarningsJSON []WarningJSON

type WarningJSON struct {
	Category string `json:"category,omitempty"`
	Tag      string `json:"tag,omitempty"`
}
