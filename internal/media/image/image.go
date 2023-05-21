package image

import (
	"database/sql/driver"
	"encoding/json"
)

type Sizes struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Image struct {
	Id          string `json:"id"`
	FileName    string `json:"fileName"`
	Sizes       Sizes  `json:"sizes"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

type NullImage struct {
	Image Image
	Valid bool
}

func (i *NullImage) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	return json.Unmarshal(value.([]byte), i)
}

func (i NullImage) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}

	return json.Marshal(i)
}

func (i *NullImage) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	i.Valid = true
	return json.Unmarshal(data, &i.Image)
}

func (i NullImage) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Image)
	}

	return []byte("null"), nil
}

type Images []Image

func (i *Images) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	return json.Unmarshal(value.([]byte), i)
}

func (i Images) Value() (driver.Value, error) {
	return json.Marshal(i)
}
