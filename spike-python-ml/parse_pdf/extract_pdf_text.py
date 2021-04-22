import pdfplumber
import sys
import parse_pdf.hayd.context as hayd_context
import parse_pdf.csv as parse_csv

if __name__ == "__main__":

  context = hayd_context.ParseContext()

  table_settings = {
      "vertical_strategy": "lines",
      "horizontal_strategy": "lines",
      "explicit_vertical_lines": [0,300,600],
      "explicit_horizontal_lines": [0,50, 720, 775],
  }

  with pdfplumber.open(sys.argv[1]) as pdf:
    # we only care about pages 3 (0th index) onward
    idx = 3
    for page in pdf.pages[3:59]:
      print("Processing page " + str(idx))
      table = page.extract_tables(table_settings=table_settings)
      left_col = table[0][1][0]
      left_col_lines = left_col.split("\n")
      for line in left_col_lines:
        context.process_raw_line(line)


      right_col = table[0][1][1]
      right_col_lines = right_col.split("\n")
      for line in right_col_lines:
        context.process_raw_line(line)

      idx += 1
    
  # now write to the output file
  parse_csv.write_table_as_csv(parse_csv.convert_to_table(context.get_parsed_catalog()))
