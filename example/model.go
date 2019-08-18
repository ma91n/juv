//go:generate juv $GOFILE
package example

import (
	"time"
)

type Post struct {
	ID      int       `validate:"gt=0" json:"id"`
	Created time.Time `validate:"max=11" json:"created"`
	Draft   bool      `json:"draft"`
	Title   string    `validate:"required"`
	Body    string    `validate:"max=50"`
}

//func (r *Post) UnmarshalJSON(b []byte) error {
//	// ------------ ここからは定形 ---------------
//	type Pointer Post // Avoid stack over flow
//	var p Pointer
//
//	if err := json.Unmarshal(b, &p); err != nil {
//		return err
//	}
//	// ------------ 定形 -------------------------
//
//	if err := validator.New().Struct(p); err != nil {
//		return err
//	}
//
//	// レシーバー変数に代入する処理
//	r.ID = p.ID
//	r.Created = p.Created
//	r.Draft = p.Draft
//	r.Body = p.Body
//
//	return nil
//}
