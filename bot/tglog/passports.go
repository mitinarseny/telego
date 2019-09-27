package tglog

import "time"

type PassportElementType string

const (
    PersonalDetailsPassportElementType       PassportElementType = "personal_details"
    PassportPassportElementType              PassportElementType = "passport"
    DriverLicensePassportElementType         PassportElementType = "driver_license"
    IdentityCardPassportElementType          PassportElementType = "identity_card"
    InternalPassportPassportElementType      PassportElementType = "internal_passport"
    AddressPassportElementType               PassportElementType = "address"
    UtilityBillPassportElementType           PassportElementType = "utility_bill"
    BankStatementPassportElementType         PassportElementType = "bank_statement"
    RentalAgreementPassportElementType       PassportElementType = "rental_agreement"
    PassportRegistrationPassportElementType  PassportElementType = "passport_registration"
    TemporaryRegistrationPassportElementType PassportElementType = "temporary_registration"
    PhoneNumberPassportElementType           PassportElementType = "phone_number"
    EmailPassportElementType                 PassportElementType = "email"
)

type PassportData struct {
    Data        *EncryptedPassportElement `bson:"data,omitempty"`
    Credentials *EncryptedCredentials     `bson:"credentials,omitempty"`
}

type EncryptedPassportElement struct {
    Type        PassportElementType `bson:"type,omitempty"`
    Data        *string             `bson:"data,omitempty"`
    PhoneNumber *string             `bson:"phone_number,omitempty"`
    Email       *string             `bson:"email,omitempty"`
    Files       []PassportFile      `bson:"files,omitempty"`
    FrontSide   []PassportFile      `bson:"front_side,omitempty"`
    ReverseSide []PassportFile      `bson:"reverse_side,omitempty"`
    Selfie      []PassportFile      `bson:"selfie,omitempty"`
    Translation []PassportFile      `bson:"translation,omitempty"`
    Hash        string              `bson:"hash,omitempty"`
}

type PassportFile struct {
    FileID   string    `bson:"file_id,omitempty"`
    FileSize int64     `bson:"file_size,omitempty"`
    FileDate time.Time `bson:"file_date,omitempty"`
}

type EncryptedCredentials struct {
    Data   string `bson:"data,omitempty"`
    Hash   string `bson:"hash,omitempty"`
    Secret string `bson:"secret,omitempty"`
}
