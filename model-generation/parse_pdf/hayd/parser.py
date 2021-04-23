import regex
from typing import Dict, List

# define lines we know to skip
KNOWN_SKIP_LINES = set(
    [
        "Product Selection",
        # the next are the second lines of long sub categories,
        # we'll use hacks to address this
        "Butterscotch & Coconut)",
        "Spread",
    ]
)


def handle_no_space_before_sub_cat(phrase, break_idx):
    def _func(line):
        phrase_idx = line.find(phrase)
        left = line[: phrase_idx + break_idx]
        right = line[phrase_idx + break_idx :]
        return (left, right)

    return _func


def strip_unnecessary_bullet(phrase):
    return (None, BULLET_REGEX.sub("", phrase))


# a hack to override the sub cat names when they are multiline in length
OVERRIDE_SUB_CAT = {
    "Baking Nuts & Chips (Chocolate,": "Baking Nuts & Chips",
    "Mayonnaise & Miracle Whip & Sandwich": "Sandwich Spread",
    "Snack Cakes (Hostess & Little Debbie)": "Snack Cakes",
}

# a list of bandages to handle weird one off cases
NEW_SUB_CATEGORY_HACKS = [
    {
        "regex": regex.compile("ozBaking"),
        "lambda": handle_no_space_before_sub_cat("ozBaking", 2),
    },
    {"regex": regex.compile("o  Fish – Fresh"), "lambda": strip_unnecessary_bullet},
    {"regex": regex.compile("o  Cream Cheese"), "lambda": strip_unnecessary_bullet},
]

# a dict of names that will conflict with the bullet regex
ITEM_HACKS = tuple([
    "Arg",
    "Nach",
    "Frit",
    "Jalapen",
    "Kashi G",
    "Ore",
    "Nabisc",
    "Old El Pas",
    "Ready T",
    "Potat",
    "Preg",
    "Progress",
    "Sargent",
])


def apply_item_fix(tokens: List[str]):
    retval = []

    idx = 0
    while idx < len(tokens):
        token = tokens[idx]

        if token.endswith(ITEM_HACKS) and idx < len(tokens) - 1:
            retval.append(f"{token}o {tokens[idx+1]}")
            idx += 1
        else:
            retval.append(token)

        idx += 1

    return retval


BULLET_REGEX = regex.compile("o  ?")

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


class Parser:
    def __init__(self):
        self.parsed_data: Dict[str, Dict[str, List[str]]] = {}
        self.current_category: str = None
        self.current_sub_category: str = None
        self.current_item: str = None

    def _finalize_current_item(self):

        # santiy check: if the current item is none, just return but warn that we have a logic error
        if self.current_item is None:
            # print("Null current item was almost finalized, skipping")
            return

        # sanity check: make sure we have a cat and sub cat
        if self.current_category is None:
            raise Exception("attempted to finalize current items without cat")

        if self.current_sub_category is None:
            raise Exception("attempted to finalize current items without sub cat")

        if self.current_category not in self.parsed_data:
            self.parsed_data[self.current_category] = {}

        if self.current_sub_category not in self.parsed_data[self.current_category]:
            self.parsed_data[self.current_category][self.current_sub_category] = []

        self.parsed_data[self.current_category][self.current_sub_category].append(
            self.current_item.strip()
        )
        self.current_item = None

    def _update_current_item(self, data: str):
        if self.current_item is None:
            self.current_item = data

        else:
            self.current_item = f"{self.current_item} {data}"

    def process_raw_line(self, raw_line: str):

        line = raw_line.strip()

        # skip known words
        if line in KNOWN_SKIP_LINES:
            return

        # if it matches a known category set it
        if line in KNOWN_CATEGORIES:
            self._finalize_current_item()
            self.current_category = line
            self.current_sub_category = None
            return

        # if it's a blank line, most like a new sub category
        if line == "":
            self._finalize_current_item()
            self.current_sub_category = None
            return

        if self.current_sub_category is None and not BULLET_REGEX.match(line):
            new_sub_cat = (
                line if line not in OVERRIDE_SUB_CAT else OVERRIDE_SUB_CAT[line]
            )
            self.current_sub_category = new_sub_cat
            return

        # if there's a weird case where it didn't space correctly
        hacks = [b for b in NEW_SUB_CATEGORY_HACKS if b["regex"].match(line)]

        if len(hacks) > 1:
            raise Exception("can't handle multiple hacks on same line")

        if len(hacks) == 1:
            (left, right) = hacks[0]["lambda"](line)
            if left is not None:
                self._update_current_item(left)
            self._finalize_current_item()
            self.current_sub_category = right
            return

        # test to see if the line has one or more bullets
        bullets: List[str] = [m.start() for m in BULLET_REGEX.finditer(line)]

        # if there's no bullets
        if len(bullets) == 0:
            self._update_current_item(line)
            return

        # if the bullet is at the beginning, process new line
        line_starts_new_bullet = bullets[0] == 0

        # split the line by bullets
        # this should now have something like ["item a", "item b"]
        # or ["cnt, 1 each", "Store brand English muffins"]
        temp_items = list(
            filter(
                None,
                BULLET_REGEX.split(
                    line,
                    0,
                ),
            )
        )

        temp_items = apply_item_fix(temp_items)

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
