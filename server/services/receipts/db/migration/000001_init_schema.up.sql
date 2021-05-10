CREATE TABLE IF NOT EXISTS unparsed_receipt_requests (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL,
  original_url VARCHAR(255) NOT NULL,
  request_timestamp TIMESTAMPTZ NOT NULL,
  raw_html TEXT NOT NULL,
  UNIQUE(original_url)
);

CREATE TABLE IF NOT EXISTS parsed_receipts (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  order_number VARCHAR(255) NOT NULL,
  order_timestamp TIMESTAMPTZ NOT NULL,
  sales_tax NUMERIC(10,2) NOT NULL,
  tip NUMERIC(10,2) NOT NULL,
  service_fee NUMERIC(10,2) NOT NULL,
  delivery_fee NUMERIC(10,2) NOT NULL,
  discounts NUMERIC(10,2) NOT NULL,
  unparsed_receipt_request_id uuid REFERENCES unparsed_receipt_requests(id) NOT NULL,
  UNIQUE(order_number)
);

CREATE TABLE IF NOT EXISTS parsed_items (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  total_cost numeric(10,2) NOT NULL,
  parsed_receipt_id uuid REFERENCES parsed_receipts(id) NOT NULL,
  category VARCHAR(255),
  unit_cost NUMERIC(10,2),
  qty INT,
  weight numeric(10,2),
  container_size NUMERIC(10,2),
  container_unit VARCHAR(25)
);
