package repo

import "time"

type BaseModel struct {
    CreatedAt time.Time `bson:"createdAt,omitempty"`
    UpdatedAt time.Time `bson:"updatesAt,omitempty"`
}
