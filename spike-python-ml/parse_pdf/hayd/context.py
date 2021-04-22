import re

# define lines we know to skip
KNOWN_SKIP_WORDS = set([
        "Product Selection"
    ]
)

def handle_no_space_before_sub_cat(phrase, breakIdx):
    def _func(line):
        phraseIdx = line.find(phrase)
        left = line[:phraseIdx+breakIdx]
        right = line[phraseIdx + breakIdx:]
        return (left, right)
    return _func

def strip_unnecessary_bullet(phrase):
    return (None, BULLET_REGEX.sub('', phrase))

# a list of bandages to handle weird one off cases
KNOWN_NEW_SUB_CATEGORY_BREAKS = [
    {
        "regex": re.compile("ozBaking"),
        "lambda": handle_no_space_before_sub_cat("ozBaking", 2)
    },
    {
        "regex": re.compile("o  Fish – Fresh"),
        "lambda": strip_unnecessary_bullet
    },
    
]

BULLET_REGEX = re.compile("o  ?")

# define known categories based on table of contents
KNOWN_CATEGORIES = set(
    [
        "Bakery",
        "Baking & Cooking",
        "Beverages",
        "Breads, Buns, Bagels",
        "Cereal & Breakfast Foods",
        "Coffee, Tea, Cocoa",
        "Condiments",
        "Deli",
        "Ethnic Foods",
        "Fish & Seafood",
        "Frozen Foods",
        "Fruits – Fresh",
        "Fruits & Vegetables – Canned, Dry",
        "Health Care",
        "Household",
        "Hygiene",
        "Jam, Jelly, Peanut Butter & Honey",
        "Meals & Sides - Boxed, Canned",
        "Meats",
        "Milk & Dairy",
        "Pasta, Rice, Sauce",
        "Refrigerated Potatoes, Sides",
        "Requests",
        "Snacks & Desserts",
        "Soups",
        "Vegetables – Fresh",
    ]
)


class ParseContext:
    def __init__(self):
        self.parsed_data = {}
        self.currentCategory = None
        self.currentSubCategory = None
        self.currentItem = None

    def _finalize_current_item(self):

        # santiy check: if the current item is none, just return but warn that we have a logic error
        if self.currentItem is None:
            # print("Null current item was almost finalized, skipping")
            return

        # sanity check: make sure we have a cat and sub cat
        if self.currentCategory is None:
            raise Exception("attempted to finalize current items without cat")

        if self.currentSubCategory is None:
            raise Exception("attempted to finalize current items without sub cat")

        if self.currentCategory not in self.parsed_data:
            self.parsed_data[self.currentCategory] = {}

        if self.currentSubCategory not in self.parsed_data[self.currentCategory]:
            self.parsed_data[self.currentCategory][self.currentSubCategory] = []

        self.parsed_data[self.currentCategory][self.currentSubCategory].append(self.currentItem)
        self.currentItem = None

    def _update_current_item(self, data):
        if self.currentItem is None:
            self.currentItem = data

        else:
            self.currentItem = f"{self.currentItem} {data}"

    def process_raw_line(self, raw_line):

        line = raw_line.strip()

        # skip known words
        if line in KNOWN_SKIP_WORDS:
            return

        # if it matches a known category set it
        elif line in KNOWN_CATEGORIES:
            print(f"** New category: {line}")
            self._finalize_current_item()
            self.currentCategory = line
            self.currentSubCategory = None
            return

        elif self.currentSubCategory is None and not BULLET_REGEX.match(line):
            print(f"*** New sub category: {line}")
            self.currentSubCategory = line
            return
        
        # if it's a blank line, most like a new sub category
        elif line == "":
            self._finalize_current_item()
            self.currentSubCategory = None
            return


        # if there's a weird case where it didn't space correctly
        breaks = [ b for b in KNOWN_NEW_SUB_CATEGORY_BREAKS if b["regex"].match(line)]

        if len(breaks) > 1:
            raise Exception("can't handle multiple break conditions")

        elif len(breaks) == 1:
            (left, right) = breaks[0]["lambda"](line)
            if left is not None:
                self._update_current_item(left)
            self._finalize_current_item()
            self.currentSubCategory = right
            return

        # test to see if the line has one or more bullets
        bullets = [m.start() for m in BULLET_REGEX.finditer(line)]

        # if there's no bullets
        if len(bullets) == 0:
            self._update_current_item(line)
            return


        # if the bullet is at the beginning, process new line
        processed_first_item = False
        line_starts_new_bullet = bullets[0] == 0

        # split the line by bullets
        # this should now have something like ["item a", "item b"]
        # or ["cnt, 1 each", "Store brand English muffins"]
        temp_items = list(filter(None, BULLET_REGEX.split(line, 0, )))

        # if we have a situation like "end"<bullet>"begin"...
        # finish the last part the item, finalize and start the new one
        start_idx = 0
        if not line_starts_new_bullet and len(temp_items) >= 2:
            self._update_current_item(temp_items[0])
            self._finalize_current_item()
            self._update_current_item(temp_items[1])
            start_idx = 2
        
        for temp_item in temp_items[start_idx:]:
            self._finalize_current_item()
            self._update_current_item(temp_item)


    def get_parsed_catalog(self):
        return self.parsed_data