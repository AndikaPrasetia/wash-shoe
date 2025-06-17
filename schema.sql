-------------------------------
-- 0. Buat Schema Auth
-------------------------------
CREATE SCHEMA IF NOT EXISTS auth;

-- Tabel utama untuk user authentication
CREATE TABLE IF NOT EXISTS auth.users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  confirmed_at TIMESTAMPTZ,
  last_sign_in_at TIMESTAMPTZ
);

-------------------------------
-- 1. Enable Extensions
-------------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-------------------------------
-- 2. Enums
-------------------------------
CREATE TYPE order_status AS ENUM (
  'pending', 'processing', 'cleaning', 'ready_for_delivery', 'completed', 'delivered', 'cancelled'
);

CREATE TYPE payment_status AS ENUM ('pending', 'success', 'failed', 'refunded');
CREATE TYPE payment_method AS ENUM ('DANA', 'OVO', 'COD', 'Bank_Transfer', 'QRIS');
CREATE TYPE notification_channel AS ENUM ('email', 'whatsapp', 'in_app');
CREATE TYPE photo_type AS ENUM ('before', 'after');

-------------------------------
-- 3. Application Tables
-------------------------------
CREATE TABLE public.users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  full_name TEXT,
  email TEXT UNIQUE NOT NULL,
  phone_number TEXT,
  provider TEXT DEFAULT 'email',
  provider_id TEXT,
  role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_auth_user FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- Addresses
CREATE TABLE public.addresses (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  street TEXT NOT NULL,
  city TEXT NOT NULL,
  province TEXT NOT NULL,
  postal_code TEXT,
  notes TEXT,
  is_primary BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Service Types
CREATE TABLE public.service_types (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT,
  base_price NUMERIC(10,2) NOT NULL,
  estimated_duration_hours INT,
  is_eco_friendly BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Orders
CREATE TABLE public.orders (
  id SERIAL PRIMARY KEY,
  user_id UUID REFERENCES public.users(id) ON DELETE SET NULL,
  address_id INT REFERENCES public.addresses(id) ON DELETE SET NULL,
  total_price NUMERIC(10,2) NOT NULL,
  status order_status DEFAULT 'pending',
  is_express BOOLEAN DEFAULT false,
  express_fee NUMERIC(10,2) DEFAULT 0,
  promo_code TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  completed_at TIMESTAMPTZ
);

-- Order Services
CREATE TABLE public.order_services (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL REFERENCES public.orders(id) ON DELETE CASCADE,
  service_type_id INT NOT NULL REFERENCES public.service_types(id) ON DELETE CASCADE,
  quantity INT DEFAULT 1,
  price NUMERIC(10,2) NOT NULL
);

-- Payments
CREATE TABLE public.payments (
  id SERIAL PRIMARY KEY,
  order_id INT UNIQUE NOT NULL REFERENCES public.orders(id) ON DELETE CASCADE,
  method payment_method NOT NULL,
  amount NUMERIC(10,2) NOT NULL,
  transaction_id TEXT,
  status payment_status DEFAULT 'pending',
  paid_at TIMESTAMPTZ
);

-- Promos
CREATE TABLE public.promos (
  code TEXT PRIMARY KEY,
  discount_type TEXT NOT NULL CHECK (discount_type IN ('percentage', 'fixed')),
  discount_value NUMERIC(10,2) NOT NULL,
  max_usage INT,
  used_count INT DEFAULT 0,
  valid_until TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Reviews
CREATE TABLE public.reviews (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  order_id INT UNIQUE NOT NULL REFERENCES public.orders(id) ON DELETE CASCADE,
  rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
  comment TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Notifications
CREATE TABLE public.notifications (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  order_id INT REFERENCES public.orders(id) ON DELETE CASCADE,
  message TEXT NOT NULL,
  channel notification_channel NOT NULL,
  status TEXT DEFAULT 'pending',
  sent_at TIMESTAMPTZ
);

-- Photo Evidences
CREATE TABLE public.photo_evidences (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL REFERENCES public.orders(id) ON DELETE CASCADE,
  photo_url TEXT NOT NULL,
  type photo_type NOT NULL,
  uploaded_by UUID REFERENCES public.users(id) ON DELETE SET NULL,
  uploaded_at TIMESTAMPTZ DEFAULT NOW()
);

-- Order Status History
CREATE TABLE public.order_status_history (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL REFERENCES public.orders(id) ON DELETE CASCADE,
  status order_status NOT NULL,
  updated_by UUID REFERENCES public.users(id) ON DELETE SET NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Blog Posts
CREATE TABLE public.blog_posts (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  author_id UUID REFERENCES public.users(id) ON DELETE SET NULL,
  published_at TIMESTAMPTZ DEFAULT NOW(),
  is_published BOOLEAN DEFAULT false
);

-- Complaints
CREATE TABLE public.complaints (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  order_id INT NOT NULL REFERENCES public.orders(id) ON DELETE CASCADE,
  description TEXT NOT NULL,
  status TEXT DEFAULT 'open' CHECK (status IN ('open', 'resolved', 'closed')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  resolved_at TIMESTAMPTZ
);

-- Referrals
CREATE TABLE public.referrals (
  id SERIAL PRIMARY KEY,
  referrer_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  referred_email TEXT NOT NULL,
  is_completed BOOLEAN DEFAULT false,
  cashback_amount NUMERIC(10,2) DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-------------------------------
-- 4. Refresh Tokens & Audit Log
-------------------------------
CREATE TABLE auth.refresh_tokens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  revoked BOOLEAN DEFAULT false
);

CREATE TABLE auth.audit_log (
  id SERIAL PRIMARY KEY,
  actor_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
  action TEXT NOT NULL,
  details JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-------------------------------
-- 5. Cascade Delete Function
-------------------------------
CREATE OR REPLACE FUNCTION delete_user_cascade()
RETURNS TRIGGER AS $$
BEGIN
  DELETE FROM public.users WHERE id = OLD.id;
  DELETE FROM public.addresses WHERE user_id = OLD.id;
  DELETE FROM public.orders WHERE user_id = OLD.id;
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_auth_user_delete
BEFORE DELETE ON auth.users
FOR EACH ROW EXECUTE FUNCTION delete_user_cascade();

-------------------------------
-- 6. Indexes (Critical for Performance)
-------------------------------
-- Users
CREATE INDEX idx_users_email ON public.users (email);

-- Orders
CREATE INDEX idx_orders_user_status ON public.orders (user_id, status);
CREATE INDEX idx_orders_created ON public.orders (created_at);

-- Payments
CREATE INDEX idx_payments_status ON public.payments (status);
CREATE INDEX idx_payments_transaction ON public.payments (transaction_id);

-- Order Services
CREATE INDEX idx_order_services_order ON public.order_services (order_id);

-- Addresses
CREATE INDEX idx_addresses_user ON public.addresses (user_id);

-- Reviews
CREATE INDEX idx_reviews_user ON public.reviews (user_id);
CREATE INDEX idx_reviews_order ON public.reviews (order_id);

-- Notifications
CREATE INDEX idx_notifications_user ON public.notifications (user_id);

-- Photo Evidences
CREATE INDEX idx_photo_evidences_order ON public.photo_evidences (order_id);

-- Order Status History
CREATE INDEX idx_order_status_history_order ON public.order_status_history (order_id);

-- Blog Posts
CREATE INDEX idx_blog_posts_author ON public.blog_posts (author_id);
CREATE INDEX idx_blog_posts_published ON public.blog_posts (published_at);

-- Complaints
CREATE INDEX idx_complaints_user ON public.complaints (user_id);
CREATE INDEX idx_complaints_order ON public.complaints (order_id);

-- Referrals
CREATE INDEX idx_referrals_referrer ON public.referrals (referrer_id);
CREATE INDEX idx_referrals_email ON public.referrals (referred_email);

-- Refresh Tokens
CREATE INDEX idx_refresh_token_user ON auth.refresh_tokens (user_id);
CREATE INDEX idx_refresh_token_expiry ON auth.refresh_tokens (expires_at);

-------------------------------
-- 7. Business Transaction Functions
-------------------------------
-- Fungsi untuk membuat order baru (Atomic operation)
CREATE OR REPLACE FUNCTION create_order_with_services(
  p_user_id UUID,
  p_address_id INT,
  p_total_price NUMERIC(10,2),
  p_status order_status,
  p_is_express BOOLEAN,
  p_express_fee NUMERIC(10,2),
  p_promo_code TEXT,
  p_services JSONB -- Format: [{"service_type_id": 1, "quantity": 2, "price": 50000}, ...]
) RETURNS INT AS $$
DECLARE
  new_order_id INT;
  service_item JSONB;
BEGIN
  -- Mulai transaksi implisit
  BEGIN
    -- Insert order utama
    INSERT INTO public.orders (
      user_id, address_id, total_price, status, 
      is_express, express_fee, promo_code
    ) VALUES (
      p_user_id, p_address_id, p_total_price, p_status,
      p_is_express, p_express_fee, p_promo_code
    ) RETURNING id INTO new_order_id;

    -- Insert order services dari JSONB
    FOR service_item IN SELECT * FROM jsonb_array_elements(p_services)
    LOOP
      INSERT INTO public.order_services (
        order_id, service_type_id, quantity, price
      ) VALUES (
        new_order_id,
        (service_item->>'service_type_id')::INT,
        (service_item->>'quantity')::INT,
        (service_item->>'price')::NUMERIC
      );
    END LOOP;

    -- Update promo usage jika ada promo
    IF p_promo_code IS NOT NULL THEN
      UPDATE public.promos
      SET used_count = used_count + 1
      WHERE code = p_promo_code;
    END IF;

    -- Return order ID jika sukses
    RETURN new_order_id;
    
  EXCEPTION
    WHEN others THEN
      -- Rollback otomatis jika ada error
      RAISE EXCEPTION 'Order creation failed: %', SQLERRM;
  END;
END;
$$ LANGUAGE plpgsql;

-- Fungsi untuk update status order dengan history tracking
CREATE OR REPLACE FUNCTION update_order_status(
  p_order_id INT,
  p_new_status order_status,
  p_updated_by UUID
) RETURNS VOID AS $$
BEGIN
  -- Mulai transaksi implisit
  BEGIN
    -- Update status utama
    UPDATE public.orders
    SET status = p_new_status
    WHERE id = p_order_id;

    -- Tambahkan history
    INSERT INTO public.order_status_history (order_id, status, updated_by)
    VALUES (p_order_id, p_new_status, p_updated_by);
    
  EXCEPTION
    WHEN others THEN
      RAISE EXCEPTION 'Status update failed: %', SQLERRM;
  END;
END;
$$ LANGUAGE plpgsql;

-- Fungsi untuk proses pembayaran lengkap
CREATE OR REPLACE FUNCTION process_payment(
  p_order_id INT,
  p_method payment_method,
  p_amount NUMERIC(10,2),
  p_transaction_id TEXT,
  p_status payment_status
) RETURNS VOID AS $$
BEGIN
  -- Mulai transaksi implisit
  BEGIN
    -- Insert payment record
    INSERT INTO public.payments (
      order_id, method, amount, transaction_id, status, paid_at
    ) VALUES (
      p_order_id, p_method, p_amount, p_transaction_id, p_status, NOW()
    );

    -- Update order status jika pembayaran sukses
    IF p_status = 'success' THEN
      PERFORM update_order_status(p_order_id, 'processing', NULL);
    END IF;
    
  EXCEPTION
    WHEN others THEN
      RAISE EXCEPTION 'Payment processing failed: %', SQLERRM;
  END;
END;
$$ LANGUAGE plpgsql;

-------------------------------
-- 8. Security & Realtime
-------------------------------
ALTER TABLE public.orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY "Users can manage their own orders" ON public.orders
FOR ALL USING (auth.uid() = user_id);

ALTER PUBLICATION supabase_realtime ADD TABLE public.orders;
ALTER PUBLICATION supabase_realtime ADD TABLE public.order_status_history;
