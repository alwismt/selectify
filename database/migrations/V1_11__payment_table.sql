CREATE TABLE payments (
      id BIGSERIAL PRIMARY KEY,

      order_id BIGINT NOT NULL
          REFERENCES orders(id),

      provider VARCHAR(30) NOT NULL,
      provider_payment_id VARCHAR(255) NOT NULL,
      client_secret VARCHAR(255) NULL,

      status VARCHAR(50) NOT NULL DEFAULT 'pending',

      amount BIGINT NOT NULL,
      currency CHAR(3) NOT NULL,

      failure_code VARCHAR(100),
      failure_message TEXT,

      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      paid_at TIMESTAMPTZ,

      UNIQUE (provider, provider_payment_id)
);