""" Extract products from stop 'n shop's site map """

import gzip
import xml.etree.ElementTree as ET
import re
from typing import List, Tuple

import requests

from training.categorize.fetch.export import write_table_as_csv


SNS_CAT_TO_CANONICAL_CAT_ID = {
    "frozen": 1,
    "deli": 2,
    "health-beauty": 3,
    "baking-cooking-needs": 4,
    "beverages": 5,
    "rice-grains-pasta-beans": 6,
    "condiments-sauces": 7,
    "bread-bakery": 8,
    "snacks-candy": 9,
    "pet-store": 10,
    "breakfast-cereal": 11,
    "laundry-paper-cleaning": 12,
    "dairy": 13,
    "soups-canned-goods": 14,
    "alcoholic-beverages": 15,
    "home-office": 16,
    "produce": 17,
    "baby-childcare": 18,
    "seafood": 19,
    "meat": 20,
    "meal-kits": 21,
    "floral-garden": 22,
    "charitable-contributions": 23,
}

PRODUCT_SITE_MAP_URL = "https://stopandshop.com/groceries/products-sitemap.xml.gz"


def get_site_map() -> ET.Element:
    resp = requests.get(PRODUCT_SITE_MAP_URL).content
    xml_content = gzip.decompress(resp).decode("utf-8")

    # remove blank namespaces
    xml_content = re.sub(r'\sxmlns="[^"]+"', "", xml_content, count=1)

    root = ET.fromstring(xml_content)
    return root


def process_child(e: ET.Element) -> Tuple[str, int, str]:

    loc = e.find("loc")

    if loc is None:
        return

    loc_split = loc.text.replace("https://stopandshop.com/", "").split("/")
    cat = loc_split[1]
    canonical_cat_id = SNS_CAT_TO_CANONICAL_CAT_ID[cat] if cat in SNS_CAT_TO_CANONICAL_CAT_ID else None
    if canonical_cat_id is None:
        raise ValueError(f"No canonical cat ID for {cat}")

    sub_cat = loc_split[2]
    item = loc_split[-1].replace(".html", "").replace("-", " ")

    return (item, canonical_cat_id, sub_cat)


def run(output_path: str):

    root = get_site_map()

    output = []
    errors: List[ValueError] = []
    for child in root:
        try:
            output.append(process_child(child))
        except ValueError as e:
            errors.append(e)

    write_table_as_csv(output, f"{output_path}/sns.csv")

    if len(errors) > 0:
        for err in errors:
            print(str(err))
        raise RuntimeError("Missing mappings, please review logs and update mappings")
