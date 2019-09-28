package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Invoice struct {
    Title          string `bson:"title,omitempty"`
    Description    string `bson:"description,omitempty"`
    StartParameter string `bson:"start_parameter,omitempty"`
    Currency       string `bson:"currency,omitempty"`
    TotalAmount    int    `bson:"total_amount,omitempty"`
}

type PreCheckoutQuery struct {
    ID               string     `bson:"id,omitempty"`
    From             User       `bson:"from,omitempty"`
    Currency         string     `bson:"currency,omitempty"`
    TotalAmount      int        `bson:"total_amount,omitempty"`
    InvoicePayload   string     `bson:"invoice_payload,omitempty"`
    ShippingOptionID *string    `bson:"shipping_option_id,omitempty"`
    OrderInfo        *OrderInfo `bson:"order_info,omitempty"`
}

type ShippingAddress struct {
    CountryCode string `bson:"country_code,omitempty"`
    State       string `bson:"state,omitempty"`
    City        string `bson:"city,omitempty"`
    StreetLine1 string `bson:"street_line1,omitempty"`
    StreetLine2 string `bson:"street_line2,omitempty"`
    PostCode    string `bson:"post_code,omitempty"`
}

type ShippingQuery struct {
    ID              string          `bson:"id,omitempty"`
    From            User            `bson:"from,omitempty"`
    InvoicePayload  string          `bson:"invoice_payload,omitempty"`
    ShippingAddress ShippingAddress `bson:"shipping_address,omitempty"`
}

type SuccessfulPayment struct {
    Currency                string     `bson:"currency,omitempty"`
    TotalAmount             int        `bson:"total_amount,omitempty"`
    InvoicePayload          string     `bson:"invoice_payload,omitempty"`
    ShippingOptionID        *string    `bson:"shipping_option_id,omitempty"`
    OrderInfo               *OrderInfo `bson:"order_info,omitempty"`
    TelegramPaymentChargeID string     `bson:"telegram_payment_charge_id,omitempty"`
    ProviderPaymentChargeID string     `bson:"provider_payment_charge_id,omitempty"`
}

type OrderInfo struct {
    Name            *string          `bson:"name,omitempty"`
    PhoneNumber     *string          `bson:"phone_number,omitempty"`
    Email           *string          `bson:"email,omitempty"`
    ShippingAddress *ShippingAddress `bson:"shipping_address,omitempty"`
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
