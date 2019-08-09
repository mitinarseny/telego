package tg_types

type OrderInfo struct {
    Name            *string
    PhoneNumber     *string
    Email           *string
    ShippingAddress *ShippingAddress
}
