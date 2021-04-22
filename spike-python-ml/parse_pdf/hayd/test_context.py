import parse_pdf.hayd.context as context

TEST_DATA = """Product Selection 
Bakery 
Bread 
o  Store Brand Chunky Cinnamon Bread, 1 
each 
o  Store Brand French Bread, 16 oz 
o  Store Brand Rye Bread, 1 each 
o  Store Brand Vienna Bread, 16 oz 
o  Store Brand Wheat Bread, 16 oz 
o Store Brand White Bread, 1 each 
 
Cookies, Cakes & Brownies 
o  Store Brand Shells 4 cnt Dessert, 
3.25 oz 
o  Store Brand Angel Food Cake, 13 oz 
o  Store Brand Brownies, 1 each 
o  Store Brand Cookies Chocolate Chip 12 
cnt, 1 eacho  Store Brand English Toffee Cookies, 12 
cnt 
o  Store Brand M&M Cookies, 12 each 
o  Store Brand Peanut Butter Cookies, 12 
each 
 
Croissants 
o  Store Brand Croissant Mini, 8 cnt 
o  Store Brand Croissants 4 cnt, 1 each 

Baking & Cooking 
Baking Cocoa 
o  Baker's Semi-Sweet Chocolate Baking 
Bar, 4 oz 
o  Baker's Unsweetened Chocolate Baking 
Bar, 4 oz 
o  Store Brand Baking Cocoa 
Powder, 8 oz 
o  Hershey's Cocoa Natural Unsweetened, 8 
ozBaking Mixes 
o  Betty Crocker Brownie Mix Milk 
Chocolate Traditional, 18.4 oz 
o  Betty Crocker Brownie Mix Turtle , 16 
oz 
o  *New* - Betty Crocker Cake Mix Spice, 
15.25 oz 
"""


def test_parse_works():
    pc = context.ParseContext()

    lines = TEST_DATA.split("\n")
    for line in lines:
        pc.process_raw_line(line.strip())

    parsed_data = pc.get_parsed_catalog()

    expected_parsed_data = {
        "Bakery": {
            "Bread": [
                "Store Brand Chunky Cinnamon Bread, 1 each",
                "Store Brand French Bread, 16 oz",
                "Store Brand Rye Bread, 1 each",
                "Store Brand Vienna Bread, 16 oz",
                "Store Brand Wheat Bread, 16 oz",
                "Store Brand White Bread, 1 each",
            ],
            "Cookies, Cakes & Brownies": [
                "Store Brand Shells 4 cnt Dessert, 3.25 oz",
                "Store Brand Angel Food Cake, 13 oz",
                "Store Brand Brownies, 1 each",
                "Store Brand Cookies Chocolate Chip 12 cnt, 1 each",
                "Store Brand English Toffee Cookies, 12 cnt",
                "Store Brand M&M Cookies, 12 each",
                "Store Brand Peanut Butter Cookies, 12 each",
            ],
            "Croissants": [
                "Store Brand Croissant Mini, 8 cnt",
                "Store Brand Croissants 4 cnt, 1 each",
            ],
        },
        "Baking & Cooking": {
            "Baking Cocoa": [
                "Baker's Semi-Sweet Chocolate Baking Bar, 4 oz",
                "Baker's Unsweetened Chocolate Baking Bar, 4 oz",
                "Store Brand Baking Cocoa Powder, 8 oz",
                "Hershey's Cocoa Natural Unsweetened, 8 oz",
            ],
            "Baking Mixes": [
                "Betty Crocker Brownie Mix Milk Chocolate Traditional, 18.4 oz",
                "Betty Crocker Brownie Mix Turtle , 16 oz",
                "*New* - Betty Crocker Cake Mix Spice, 15.25 oz",
            ],
        },
    }

    assert parsed_data == expected_parsed_data
