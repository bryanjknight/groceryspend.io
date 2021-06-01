CREATE TABLE IF NOT EXISTS request_type (
 id SERIAL PRIMARY KEY,
 label varchar(16)
);

INSERT INTO request_type (label) VALUES ('HTML');
INSERT INTO request_type (label) VALUES ('Image');

CREATE TABLE IF NOT EXISTS status_type (
 id SERIAL PRIMARY KEY,
 label varchar(16)
);
INSERT INTO status_type (label) VALUES ('Submitted');
INSERT INTO status_type (label) VALUES ('Processing');
INSERT INTO status_type (label) VALUES ('Completed');
INSERT INTO status_type (label) VALUES ('Error');

-- create new foreign keys to the request table
ALTER TABLE unparsed_receipt_requests ADD COLUMN request_type_id INT;
ALTER TABLE unparsed_receipt_requests ADD COLUMN status_type_id INT;

-- -- update request table with the request type, (they'll all be html and complete)
UPDATE unparsed_receipt_requests
  SET 
    request_type_id = 1,
    status_type_id = 3;


-- add foreign key constraints
ALTER TABLE unparsed_receipt_requests ADD CONSTRAINT fk_request_id FOREIGN KEY (request_type_id) REFERENCES request_type(id);
ALTER TABLE unparsed_receipt_requests ADD CONSTRAINT fk_status_id FOREIGN KEY (status_type_id) REFERENCES status_type(id);