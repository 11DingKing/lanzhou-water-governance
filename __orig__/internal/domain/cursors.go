package domain

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type Cursor struct {
	CreatedAt string
	ID        int64
}

func (c Cursor) Encode() string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%s|%d", c.CreatedAt, c.ID)))
}
func DecodeCursor(value string) (Cursor, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return Cursor{}, err
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 2 {
		return Cursor{}, fmt.Errorf("invalid cursor")
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	return Cursor{CreatedAt: parts[0], ID: id}, err
}
func NextCursor(createdAt string, id int64) Cursor { return Cursor{CreatedAt: createdAt, ID: id} }
func EmptyCursor(value string) bool                { return strings.TrimSpace(value) == "" }
