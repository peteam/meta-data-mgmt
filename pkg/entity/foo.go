package entity

type Foo struct {
	Attr1            string `json:"attr1"`
	Attr2            string `json:"attr2"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
	UpdatedTimestamp int64  `json:"updatedTimestamp"`
}
