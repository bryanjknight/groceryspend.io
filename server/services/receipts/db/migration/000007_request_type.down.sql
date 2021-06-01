ALTER TABLE unparsed_receipt_requests DROP COLUMN request_type_id;
ALTER TABLE unparsed_receipt_requests DROP COLUMN status_type_id;

DROP TABLE IF EXISTS request_type;
DROP TABLE IF EXISTS status_type;