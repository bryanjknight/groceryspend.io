ALTER TABLE unparsed_receipt_requests DROP CONSTRAINT fk_request_id;
ALTER TABLE unparsed_receipt_requests DROP CONSTRAINT fk_status_id;
ALTER TABLE parsed_receipts DROP CONSTRAINT parsed_receipts_unparsed_receipt_request_id_fkey;
ALTER TABLE parsed_items DROP CONSTRAINT parsed_items_parsed_receipt_id_fkey;

ALTER TABLE unparsed_receipt_requests 
  ADD CONSTRAINT unparsed_receipt_requests_original_url_key 
  UNIQUE (original_url);

ALTER TABLE parsed_receipts 
  ADD CONSTRAINT parsed_receipts_order_number_key 
  UNIQUE (order_number); 

ALTER TABLE unparsed_receipt_requests 
  ADD CONSTRAINT fk_request_id 
  FOREIGN KEY (request_type_id) 
  REFERENCES request_type(id);

ALTER TABLE unparsed_receipt_requests 
  ADD CONSTRAINT fk_status_id 
  FOREIGN KEY (status_type_id) 
  REFERENCES status_type(id);

ALTER TABLE parsed_receipts 
  ADD CONSTRAINT parsed_receipts_unparsed_receipt_request_id_fkey
  FOREIGN KEY (unparsed_receipt_request_id)
  REFERENCES unparsed_receipt_requests(ID);

ALTER TABLE parsed_items 
  ADD CONSTRAINT parsed_items_parsed_receipt_id_fkey 
  FOREIGN KEY (parsed_receipt_id)
  REFERENCES parsed_receipts(ID);