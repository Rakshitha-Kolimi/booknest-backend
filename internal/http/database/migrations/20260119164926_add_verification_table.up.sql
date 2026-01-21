CREATE TYPE VERIFICATION_TOKEN_TYPE AS ENUM (
    'EMAIL_VERIFICATION',
    'MOBILE_VERIFICATION',
    'PASSWORD_RESET',
    'LOGIN_OTP'
);

CREATE TABLE verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VERIFICATION_TOKEN_TYPE NOT NULL,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    used_at TIMESTAMPTZ NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Indexes for performance
CREATE INDEX idx_verification_tokens_user_id ON verification_tokens(user_id);

CREATE INDEX idx_verification_tokens_type ON verification_tokens(type);

CREATE INDEX idx_verification_tokens_token_hash ON verification_tokens(token_hash);

CREATE INDEX idx_verification_tokens_expires_at ON verification_tokens(expires_at);

-- Enforce one active token per user per type
CREATE UNIQUE INDEX ux_verification_tokens_active ON verification_tokens(user_id, type)
WHERE
    is_used = FALSE;
