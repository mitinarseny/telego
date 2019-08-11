package repository

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Invoice struct {
    Title          string
    Description    string
    StartParameter string
    Currency       string
    TotalAmount    int
}

type PreCheckoutQuery struct {
    ID               string
    From             User
    Currency         string
    TotalAmount      int
    InvoicePayload   string
    ShippingOptionID *string
    OrderInfo        *OrderInfo
}

func fromPreCheckoutQuery(q *tb.PreCheckoutQuery) *PreCheckoutQuery {
    cq := new(PreCheckoutQuery)
    cq.ID = q.ID
    cq.From = *fromTelebotUser(q.Sender)
    cq.Currency = q.Currency
    cq.TotalAmount = q.Total
    cq.InvoicePayload = q.Payload
    return cq
}

type ShippingAddress struct {
    CountryCode string
    State       string
    City        string
    StreetLine1 string
    StreetLine2 string
    PostCode    string
}

type ShippingQuery struct {
    ID              string
    From            User
    InvoicePayload  string
    ShippingAddress ShippingAddress
}

type SuccessfulPayment struct {
    Currency                string
    TotalAmount             int
    InvoicePayload          string
    ShippingOptionID        *string
    OrderInfo               *OrderInfo
    TelegramPaymentChargeID string
    ProviderPaymentChargeID string
}

type OrderInfo struct {
    Name            *string
    PhoneNumber     *string
    Email           *string
    ShippingAddress *ShippingAddress
}
