#!/bin/bash
set -e

clickhouse client --user logger --password password123 -n <<-EOSQL
  CREATE DATABASE IF NOT EXISTS log;

  CREATE TABLE IF NOT EXISTS log.updates
  (
      _created_at DEFAULT now(),

      update_id                      Int64,
      message_id                     Int64,
      edited_message_id              Int64,
      channel_post_message_id        Int64,
      edited_channel_post_message_id Int64,
      inline_query_id                String,
      chosen_inline_result_id        String,
      callback_query_id              String,
      shipping_query_id              String,
      pre_checkout_query_id          String,
      poll_id                        String
  ) ENGINE = MergeTree()
        PARTITION BY toYYYYMM(_created_at)
        ORDER BY update_id;

  CREATE TABLE IF NOT EXISTS log.messages
  (
      _created_at DEFAULT now(),

      message_id               Int64,
      from_user_id             Int64,
      date                     DateTime,
      chat_id                  Int64,
      forward_from_user_id     Int64,
      forward_from_chat_id     Int64,
      forward_from_message_id  Int64,
      forward_signature        String,
      forward_sender_name      String,
      forward_date             DateTime,
      reply_to_message_id      Int64,
      edit_date                DateTime,
      media_group_id           String,
      author_signature         String,
      text                     String,
      entities_ids             Array(UUID),
      caption_entities_ids     Array(UUID),
      audio_id                 UUID,
      document_id              UUID,
      animation_id             UUID,
      game_id                  UUID,
      photo_ids                Array(UUID),
      sticker_id               UUID,
      video_id                 UUID,
      voice_id                 UUID,
      video_note_id            UUID,
      caption                  String,
      contact_id               UUID,
      location_id              UUID,
      venue_id                 UUID,
      poll_id                  String,
      new_chat_members_ids     Array(Int64),
      left_chat_member_id      Int64,
      new_chat_title           String,
      new_chat_photo           Array(UUID),
      delete_chat_photo        UInt8,
      group_chat_created       UInt8,
      supergroup_chat_created  UInt8,
      channel_chat_created     UInt8,
      migrate_to_chat_id       Int64,
      migrate_from_cha_id      Int64,
      pinned_message_id        Int64,
      invoice_id               UUID,
      successful_payment_id    UUID,
      connected_website        String,
      passport_data_id         UUID,
      reply_markup_buttons_ids Array(Array(UUID))
  ) ENGINE = MergeTree()
        PARTITION BY toYYYYMM(date)
        ORDER BY message_id;

  CREATE TABLE IF NOT EXISTS log.users
  (
      _created_at DEFAULT now(),

      id            Int64,
      is_bot        UInt8,
      first_name    String,
      last_name     String,
      username      String,
      language_code FixedString(5)
  ) ENGINE MergeTree()
        ORDER BY id;

  CREATE TABLE IF NOT EXISTS log.chats
  (
      _created_at DEFAULT now(),

      id         Int64,
      type       Enum8(
          'private' = 1,
          'group' = 2,
          'supergroup' = 3,
          'channel' = 4
          ),
      title      String,
      username   String,
      first_name String,
      last_name  String
  ) ENGINE = MergeTree()
        ORDER BY id;

  CREATE TABLE IF NOT EXISTS log.message_entities
  (
      _created_at DEFAULT now(),
      _id        UUID,

      message_id Int64,
      type       Enum8(
          'cashtag' = 1,
          'bot_command' = 2,
          'url' = 3,
          'email' = 4,
          'phone_number' = 5,
          'bold' = 6,
          'italic' = 7,
          'code' = 8,
          'pre' = 9,
          'text_link' = 10,
          'text_mention' = 11
          ),
      offset     Int16,
      length     Int16,
      URL        String,
      user_id    Int64
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.audios
  (
      _created_at DEFAULT now(),
      _id       UUID,

      file_id   String,
      duration  Int16,
      performer String,
      title     String,
      mime_type String,
      file_size Int64,
      thumb_id  UUID
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.documents
  (
      _created_at DEFAULT now(),
      _id       UUID,

      file_id   String,
      thumb_id  UUID,
      file_name String,
      mime_type String,
      file_size Int64
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.animations
  (
      _created_at DEFAULT now(),
      _id       UUID,

      file_id   String,
      width     Int16,
      height    Int16,
      duration  Int16,
      thumb_id  UUID,
      file_name String,
      mime_type String,
      file_size Int64
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.games
  (
      _created_at   DateTime DEFAULT now(),
      _id           UUID,

      title         String,
      description   String,
      photo_ids     Array(UUID),
      text          String,
      text_entities Array(UUID),
      animation_id  UUID
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.stickers
  (
      _created_at DEFAULT now(),
      _id         UUID,

      file_id     String,
      width       Int16,
      height      Int16,
      is_animated UInt8,
      thumb_id    UUID,
      emoji       String,
      set_name    String,
      mask_position Nested(
          point Enum8(
              'forehead' = 1,
              'eyes' = 2,
              'mouth' = 3,
              'chin' = 4
              ),
          x_shift Float32,
          y_shift Float32,
          scale Float32
          ),
      file_size   Int64
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.videos
  (
      _created_at DEFAULT now(),
      _id       UUID,

      file_id   String,
      width     Int16,
      height    Int16,
      duration  Int16,
      thumb_id  UUID,
      mime_type String,
      file_size Int64
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.voices
  (
      _created_at DEFAULT now(),
      _id       UUID,

      file_id   String,
      duration  Int16,
      mime_type String,
      file_size Int64
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.video_notes
  (
      _created_at DEFAULT now(),
      _id       UUID,

      file_id   String,
      length    Int16,
      duration  Int16,
      thumb_id  UUID,
      file_size Int64
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.contacts
  (
      _created_at  DateTime DEFAULT now(),
      _id          UUID,

      phone_number String,
      first_name   String,
      last_name    String,
      user_id      Int64,
      vcard        String
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.locations
  (
      _created_at DEFAULT now(),
      _id       UUID,

      longitude Float32,
      latitude  Float32
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.venues
  (
      _created_at     DateTime DEFAULT now(),
      _id             UUID,

      location_id     UUID,
      title           String,
      address         String,
      foursquare_id   String,
      foursquare_type String
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.polls
  (
      _created_at DEFAULT now(),
      _id         UUID,

      id          String,
      question    String,
      options_ids Array(UUID),
      is_closed   UInt8
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.poll_options
  (
      _created_at DEFAULT now(),
      _id         UUID,

      text        String,
      voter_count Int32
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.invoices
  (
      _created_at     DateTime DEFAULT now(),
      _id             UUID,

      title           String,
      description     String,
      start_parameter String,
      currency        FixedString(3),
      total_amount    Int32
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.successful_payments
  (
      _created_at                DateTime DEFAULT now(),
      _id                        UUID,

      currency                   FixedString(3),
      total_amount               Int32,
      invoice_payload            String,
      shipping_option_id         String,
      order_info Nested(
          name String,
          phone_number String,
          email String,
          shipping_address_id UUID
          ),
      telegram_payment_charge_id String,
      provider_payment_charge_id String
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.passports_data
  (
      _created_at    DateTime DEFAULT now(),
      _id            UUID,

      data_id        UUID,
      credentials_id UUID
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.encrypted_passport_elements
  (
      _created_at  DateTime DEFAULT now(),
      _id          UUID,

      type         Enum8(
          'passport' = 1,
          'driver_license' = 2,
          'identity_card' = 3,
          'internal_passport' = 4,
          'address' = 5,
          'utility_bill' = 6,
          'bank_statement' = 7,
          'rental_agreement' = 8,
          'passport_registration' = 9,
          'temporary_registration' = 10,
          'phone_number' = 11,
          'email' = 12
          ),
      data         String,
      phone_number String,
      email        String,
      files_ids    Array(UUID),
      front_side   Array(UUID),
      reverse_side Array(UUID),
      selfie       Array(UUID),
      translation  Array(UUID),
      hash         FixedString(32)
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.passport_files
  (
      _created_at DEFAULT now(),
      _id       UUID,

      file_id   String,
      file_size Int64,
      file_date DateTime
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.encrypted_credentials
  (
      _created_at DEFAULT now(),
      _id    UUID,

      data   String,
      hash   FixedString(32),
      secret String
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.inline_keyboard_buttons
  (
      _created_at                      DateTime DEFAULT now(),
      _id                              UUID,

      text                             String,
      URL                              String,
      login_URL_id                     UUID,
      callback_data                    String,
      switch_inline_query              String,
      switch_inline_query_current_chat String,
      callback_game_id                 UUID,
      pay                              UInt8
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.login_urls
  (
      _created_at          DateTime DEFAULT now(),
      _id                  UUID,

      URL                  String,
      forward_text         String,
      bot_username         String,
      request_write_access UInt8
  ) ENGINE = MergeTree()
        ORDER BY _id;

  CREATE TABLE IF NOT EXISTS log.inline_queries
  (
      _created_at DEFAULT now(),

      id          String,
      from_id     Int64,
      location_id UUID,
      query       String,
      offset      String
  ) ENGINE = MergeTree()
        ORDER BY id;

  CREATE TABLE IF NOT EXISTS log.callback_queries
  (
      _created_at       DateTime DEFAULT now(),

      id                String,
      from_id           Int64,
      message_id        Int64,
      inline_message_id String,
      chat_instance     String,
      data              String,
      game_short_name   String
  ) ENGINE = MergeTree()
        ORDER BY id;

  CREATE TABLE IF NOT EXISTS log.shipping_queries
  (
      _created_at     DateTime DEFAULT now(),

      id              String,
      from_id         Int64,
      invoice_payload String,
      shipping_address Nested(
          country_code FixedString(2),
          state String,
          city String,
          street_line1 String,
          street_line2 String,
          post_code String
          )
  ) ENGINE = MergeTree()
        ORDER BY id;

  CREATE TABLE IF NOT EXISTS log.pre_checkout_queries
  (
      _created_at        DateTime DEFAULT now(),

      id                 String,
      from_id            Int64,
      currency           FixedString(3),
      total_amount       Int32,
      invoice_payload    String,
      shipping_option_id String,
      order_info_id      UUID
  ) ENGINE = MergeTree()
        ORDER BY id;

  CREATE TABLE IF NOT EXISTS log.orders_info
  (
      _created_at  DateTime DEFAULT now(),
      _id          UUID,

      name         String,
      phone_number String,
      email        String,
      shipping_address Nested(
          country_code FixedString(2),
          state String,
          city String,
          street_line1 String,
          street_line2 String,
          post_code String
          )
  ) ENGINE = MergeTree()
        ORDER BY _id;
EOSQL
