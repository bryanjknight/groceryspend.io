-- add total cost column
ALTER TABLE parsed_receipts ADD COLUMN IF NOT EXISTS subtotal_cost NUMERIC(10,2);

-- update the total cost basedon each parsed receipt
WITH all_subtotals AS (
  SELECT parsed_receipt_id, SUM(total_cost) as subtotal
  FROM parsed_items
  GROUP BY parsed_receipt_id
)
UPDATE parsed_receipts as pr
SET subtotal_cost = pi.subtotal
FROM all_subtotals pi
WHERE pr.id = pi.parsed_receipt_id;

-- make subtotal cost not null
ALTER TABLE parsed_receipts ALTER COLUMN subtotal_cost SET NOT NULL;