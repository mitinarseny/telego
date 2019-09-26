package tg

type Contact struct {
    PhoneNumber string  `bson:"phone_number,omitempty"`
    FirstName   string  `bson:"first_name,omitempty"`
    LastName    *string `bson:"last_name,omitempty"`
    UserID      *int64  `bson:"user_id,omitempty"`
    VCard       *string `bson:"vcard,omitempty"`
}
