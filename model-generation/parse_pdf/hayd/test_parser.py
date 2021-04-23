import parse_pdf.hayd.parser as parser
import pytest

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
    pc = parser.Parser()

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

@pytest.mark.parametrize(
    ["input", "output", "name"],
    [
        [
            ["Nabisc", "Chips Ahoy! Chewy Cookies,"],
            ["Nabisco Chips Ahoy! Chewy Cookies,"],
            "Nabisco"
        ],
        [
            ["Nabisc", "Golden Ore", "Sandwich Cookies, 14.3 oz"],
            ["Nabisco Golden Oreo Sandwich Cookies, 14.3 oz"],
            "Nabisco and Oreo"
        ],        
        [
            ["Arg", "Double Acting Aluminium Free"],
            ["Argo Double Acting Aluminium Free"],
            "Argo"
        ],
        [
            ["Ortega Jalapen", "Pepper Diced, 4 oz"],
            ["Ortega Jalapeno Pepper Diced, 4 oz"],
            "Jalapenos"
        ],
        [
            ["a","b", "c"],
            ["a","b", "c"],
            "No op scenario"
        ],
        [
            ["Potat","b", "Potat", "c"],
            ["Potato b", "Potato c"],
            "alternating"
        ],
        [
            ["Potat","b", "Potat", "c", "Potat"],
            ["Potato b", "Potato c", "Potato"],
            "alternating with one at the end"
        ],
    ]
)
def test_apply_item_fix(input, output, name):
    actual = parser.apply_item_fix(input)
    assert actual == output, name