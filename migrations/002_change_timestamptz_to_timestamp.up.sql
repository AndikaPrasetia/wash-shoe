-- 002_change_timestamptz_to_timestamp.up.sql

-- auth.users
ALTER TABLE auth.users
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC',
  ALTER COLUMN updated_at TYPE timestamp without time zone USING updated_at AT TIME ZONE 'UTC',
  ALTER COLUMN confirmed_at TYPE timestamp without time zone USING confirmed_at AT TIME ZONE 'UTC',
  ALTER COLUMN last_sign_in_at TYPE timestamp without time zone USING last_sign_in_at AT TIME ZONE 'UTC';

-- public.users
ALTER TABLE public.users
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC',
  ALTER COLUMN updated_at TYPE timestamp without time zone USING updated_at AT TIME ZONE 'UTC';

-- public.addresses
ALTER TABLE public.addresses
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC';

-- public.service_types
ALTER TABLE public.service_types
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC';

-- public.orders
ALTER TABLE public.orders
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC',
  ALTER COLUMN completed_at TYPE timestamp without time zone USING completed_at AT TIME ZONE 'UTC';

-- public.payments
ALTER TABLE public.payments
  ALTER COLUMN paid_at TYPE timestamp without time zone USING paid_at AT TIME ZONE 'UTC';

-- public.promos
ALTER TABLE public.promos
  ALTER COLUMN valid_until TYPE timestamp without time zone USING valid_until AT TIME ZONE 'UTC',
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC';

-- public.reviews
ALTER TABLE public.reviews
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC';

-- public.notifications
ALTER TABLE public.notifications
  ALTER COLUMN sent_at TYPE timestamp without time zone USING sent_at AT TIME ZONE 'UTC';

-- public.photo_evidences
ALTER TABLE public.photo_evidences
  ALTER COLUMN uploaded_at TYPE timestamp without time zone USING uploaded_at AT TIME ZONE 'UTC';

-- public.order_status_history
ALTER TABLE public.order_status_history
  ALTER COLUMN updated_at TYPE timestamp without time zone USING updated_at AT TIME ZONE 'UTC';

-- public.blog_posts
ALTER TABLE public.blog_posts
  ALTER COLUMN published_at TYPE timestamp without time zone USING published_at AT TIME ZONE 'UTC';

-- public.complaints
ALTER TABLE public.complaints
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC',
  ALTER COLUMN resolved_at TYPE timestamp without time zone USING resolved_at AT TIME ZONE 'UTC';

-- public.referrals
ALTER TABLE public.referrals
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC';

-- auth.refresh_tokens
ALTER TABLE auth.refresh_tokens
  ALTER COLUMN expires_at TYPE timestamp without time zone USING expires_at AT TIME ZONE 'UTC',
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC';

-- auth.audit_log
ALTER TABLE auth.audit_log
  ALTER COLUMN created_at TYPE timestamp without time zone USING created_at AT TIME ZONE 'UTC';
