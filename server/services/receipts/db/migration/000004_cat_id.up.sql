-- add category ID
ALTER TABLE parsed_items ADD COLUMN IF NOT EXISTS category_id INT;

-- update category_id with the appropriate values
UPDATE parsed_items SET category_id =
  CASE 
    WHEN category = 'frozen' THEN 1
    WHEN category = 'deli' THEN 2
    WHEN category = 'health-beauty' THEN 3
    WHEN category = 'baking-cooking-needs' THEN 4
    WHEN category = 'beverages' THEN 5
    WHEN category = 'rice-grains-pasta-beans' THEN 6
    WHEN category = 'condiments-sauces' THEN 7
    WHEN category = 'bread-bakery' THEN 8
    WHEN category = 'snacks-candy' THEN 9
    WHEN category = 'pet-store' THEN 10
    WHEN category = 'breakfast-cereal' THEN 11
    WHEN category = 'laundry-paper-cleaning' THEN 12
    WHEN category = 'dairy' THEN 13
    WHEN category = 'soups-canned-goods' THEN 14
    WHEN category = 'alcoholic-beverages' THEN 15
    WHEN category = 'home-office' THEN 16
    WHEN category = 'produce' THEN 17
    WHEN category = 'baby-childcare' THEN 18
    WHEN category = 'seafood' THEN 19
    WHEN category = 'meat' THEN 20
    WHEN category = 'meal-kits' THEN 21
    WHEN category = 'floral-garden' THEN 22
    ELSE -1
  END;

-- make category_id not null
ALTER TABLE parsed_items ALTER COLUMN category_id SET NOT NULL;