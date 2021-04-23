from typing import Dict, List
import csv

def convert_to_table(data: Dict[str, Dict[str, List[str]]]) -> List[List[str]]:

  retval = []

  for cat, cat_data in data.items():
    for sub_cat, sub_cat_data in cat_data.items():
      for item in sub_cat_data:
        retval.append([item, cat, sub_cat])

  return retval


def write_table_as_csv(table: List[List[str]], output_path: str) -> None:
  with open(output_path, "w+") as output_fd:
    csv_writer = csv.writer(output_fd, delimiter=",")
    csv_writer.writerow(["Item", "Category", "Subcategory"])
    for row in table:
      csv_writer.writerow(row)
