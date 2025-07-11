-- 001_init_schema.down.sql

-- 5. Drop trigger & function
DROP TRIGGER IF EXISTS before_auth_user_delete ON auth.users;
DROP FUNCTION IF EXISTS delete_user_cascade();

-- 7. Drop business functions
DROP FUNCTION IF EXISTS create_order_with_services(JSONB);
DROP FUNCTION IF EXISTS update_order_status(INT, order_status, UUID);
DROP FUNCTION IF EXISTS process_payment(INT, payment_method, NUMERIC, TEXT, payment_status);
DROP FUNCTION IF EXISTS auth.sign_up_user(TEXT, TEXT, TEXT, TEXT, TEXT);

-- 6. Drop indexes (opsional—CASCADE biasanya sudah membersihkannya)
-- (bisa di-skip jika menggunakan DROP TABLE … CASCADE)

-- 4. Drop tables in auth schema
DROP TABLE IF EXISTS auth.audit_log CASCADE;
DROP TABLE IF EXISTS auth.refresh_tokens CASCADE;

-- 3. Drop application tables
DROP TABLE IF EXISTS public.referrals CASCADE;
DROP TABLE IF EXISTS public.complaints CASCADE;
DROP TABLE IF EXISTS public.blog_posts CASCADE;
DROP TABLE IF EXISTS public.order_status_history CASCADE;
DROP TABLE IF EXISTS public.photo_evidences CASCADE;
DROP TABLE IF EXISTS public.notifications CASCADE;
DROP TABLE IF EXISTS public.reviews CASCADE;
DROP TABLE IF EXISTS public.promos CASCADE;
DROP TABLE IF EXISTS public.payments CASCADE;
DROP TABLE IF EXISTS public.order_services CASCADE;
DROP TABLE IF EXISTS public.orders CASCADE;
DROP TABLE IF EXISTS public.service_types CASCADE;
DROP TABLE IF EXISTS public.addresses CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;

-- 2. Drop types
DROP TYPE IF EXISTS photo_type;
DROP TYPE IF EXISTS notification_channel;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS order_status;

-- 1. Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";

-- 0. Drop schema auth
DROP SCHEMA IF EXISTS auth CASCADE;
