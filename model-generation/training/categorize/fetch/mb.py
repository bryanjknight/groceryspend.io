"""Module for fetching Demoulas' Market Basket data

Notes:
* GET https://www.shopmarketbasket.com/departments-rest returns all departments
* GET https://www.shopmarketbasket.com/weekly-flyer-rest returns a JSON blob of data:
  [].field_featured_products[].node...
or [].field_flyer_item[].node...
"""
