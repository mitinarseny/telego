package tg

type Venue struct {
    Location       Location `bson:"location,omitempty"`
    Title          string   `bson:"title,omitempty"`
    Address        string   `bson:"address,omitempty"`
    FourSquareID   *string  `bson:"foursquare_id,omitempty"`
    FourSquareType *string  `bson:"foursquare_type,omitempty"`
}
