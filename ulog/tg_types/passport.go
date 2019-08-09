package tg_types

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
    Data        *EncryptedPassportElement
    Credentials *EncryptedCredentials
}

type EncryptedPassportElement struct {
    Type        PassportElementType
    Data        *string
    PhoneNumber *string
    Email       *string
    Files       []PassportFile
    FrontSide   []PassportFile
    ReverseSide []PassportFile
    Selfie      []PassportFile
    Translation []PassportFile
    Hash        string
}

type PassportFile struct {
    FileID   string
    FileSize int64
    FileDate time.Time
}

type EncryptedCredentials struct {
    Data   string
    Hash   string
    Secret string
}
