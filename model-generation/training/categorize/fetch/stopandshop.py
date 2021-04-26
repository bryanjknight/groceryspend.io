import gzip
import requests
import xml.etree.ElementTree as ET
import re
from typing import List
from training.categorize.export import write_table_as_csv


PRODUCT_SITE_MAP_URL = "https://stopandshop.com/groceries/products-sitemap.xml.gz"


def get_site_map() -> ET.Element:
  resp = requests.get(PRODUCT_SITE_MAP_URL).content
  xml_content = gzip.decompress(resp).decode("utf-8")

  # remove blank namespaces
  xml_content = re.sub(r'\sxmlns="[^"]+"', '', xml_content, count=1)

  root = ET.fromstring(xml_content)
  return root

def process_child(e: ET.Element) -> List[str]:

  loc = e.find("loc")
  
  if loc is None:
    return

  loc_split = loc.text.replace("https://stopandshop.com/", "").split("/")
  cat = loc_split[1]
  sub_cat = loc_split[2]
  item = loc_split[-1].replace(".html", "").replace("-", " ")

  return [item, cat, sub_cat]

if __name__ == "__main__":

  root = get_site_map()
  
output = []
for child in root:
  output.append(process_child(child))


write_table_as_csv(output, "sns.csv") 
