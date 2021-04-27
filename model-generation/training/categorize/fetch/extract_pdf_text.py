import pdfplumber
import sys
import training.categorize.fetch.hayd.parser as hayd_parser
import training.categorize.fetch.export as export_csv

if __name__ == "__main__":

    context = hayd_parser.Parser()

    table_settings = {
        "vertical_strategy": "lines",
        "horizontal_strategy": "lines",
        "explicit_vertical_lines": [0, 300, 600],
        "explicit_horizontal_lines": [0, 50, 720, 775],
    }

    with pdfplumber.open(sys.argv[1]) as pdf, open(
        sys.argv[2].replace(".csv", ".txt"), "w+"
    ) as f:
        # we only care about pages 3 (0th index) onward
        start = 3
        idx = start
        end = 59
        for page in pdf.pages[start:end]:
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

            f.write(left_col)
            f.write(right_col)
            f.flush()

            idx += 1

    context._finalize_current_item()
    # now write to the output file
    export_csv.write_table_as_csv(
        export_csv.convert_to_table(context.get_parsed_catalog()), sys.argv[2]
    )
