CREATE TABLE merchant (
    merchant_id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,

    description TEXT,
    logo TEXT,

    country_code VARCHAR(2) NOT NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    verification_status VARCHAR(50) NOT NULL DEFAULT 'pending',

    payment_provider VARCHAR(50),
    payment_provider_account_id VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT merchant_status_check
      CHECK (status IN (
                        'pending',
                        'active',
                        'suspended',
                        'closed'
          )),

    CONSTRAINT merchant_verification_status_check
      CHECK (verification_status IN (
                                     'pending',
                                     'verified',
                                     'rejected'
          ))
);

-- Remove old generic scope columns
ALTER TABLE user_role
    DROP COLUMN IF EXISTS scope_type,
    DROP COLUMN IF EXISTS scope_id;

-- Add merchant-specific columns
ALTER TABLE user_role
    ADD COLUMN merchant_role VARCHAR(50),
    ADD COLUMN merchant_id BIGINT;

ALTER TABLE user_role
    ADD CONSTRAINT user_role_merchant_fk
        FOREIGN KEY (merchant_id)
            REFERENCES merchant(merchant_id)
            ON DELETE CASCADE;

ALTER TABLE user_role
    ADD CONSTRAINT user_role_merchant_role_check
        CHECK (
            merchant_role IS NULL
                OR merchant_role IN (
                                     'owner',
                                     'admin',
                                     'manager',
                                     'staff'
                )
            );

-- one user can only belong to one merchant
ALTER TABLE user_role
    ADD CONSTRAINT user_role_user_unique
        UNIQUE (user_id);

ALTER TABLE user_role
    ADD CONSTRAINT user_role_merchant_consistency_check
        CHECK (
            (
                role = 'merchant'
                    AND merchant_id IS NOT NULL
                    AND merchant_role IS NOT NULL
                )
                OR
            (
                role <> 'merchant'
                    AND merchant_id IS NULL
                    AND merchant_role IS NULL
                )
            );