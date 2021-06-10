CREATE TABLE IF NOT EXISTS customers (
  id uuid PRIMARY KEY NOT NULL,
  stripe_customer_id varchar(255),
  UNIQUE(stripe_customer_id)
);

CREATE TABLE IF NOT EXISTS subscriptions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  customer_id uuid NOT NULL,
  subscription_type integer NOT NULL,
  created_at timestamptz NOT NULL,
  stripe_client_secret varchar(255) NOT NULL,
  canceled_at timestamptz
);

ALTER TABLE subscriptions 
  ADD CONSTRAINT fk_subscriptions_customer_id 
  FOREIGN KEY (customer_id) 
  REFERENCES customers(id)
  ON DELETE RESTRICT;

CREATE TABLE IF NOT EXISTS payments (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  subscription_id uuid NOT NULL,
  payment_date timestamptz NOT NULL,
  charge_cents integer NOT NULL,
  confirmation_code varchar(255) NOT NULL
);

ALTER TABLE payments 
  ADD CONSTRAINT fk_payments_subscription_id 
  FOREIGN KEY (subscription_id) 
  REFERENCES subscriptions(id)
  ON DELETE RESTRICT;