import json
import os
from typing import Dict

from flask import abort, Flask, request, jsonify

import joblib
import yaml

from training.categorize.train import CategorizeTuple

app = Flask(__name__)

# TODO: make this not global state and rather passed around as a context
#       for the flask routes
cat_tuple: CategorizeTuple = joblib.load(os.environ["CAT_MODEL"])
cat_id_to_cat_obj: Dict[str, Dict[str, str]] = {}
with open(os.environ["CAT_YAML"]) as f:
    y = yaml.safe_load(f)
    cat_id_to_cat_obj = {i["id"]: i for i in y["categories"]}


def init():
    pass


@app.route("/categorize", methods=["POST"])
def categorize():
    """
    Given a list of items, return a dictionary of the item and the 
    category we predict it is a member of
    """
    items = request.json
    if not items or len(items) == 0:
        return abort(400)

    item_features = cat_tuple.tfidf.transform(items)
    predictions = cat_tuple.model.predict(item_features)
    retval = {
        items[i]: cat_id_to_cat_obj[cat_tuple.id_to_cat[predictions[i]]]
        for i in range(len(predictions))
    }
    return jsonify(retval)


@app.route("/categories", methods=["GET"])
def all_categories():
    """
    Return array of all categories
    """
    return jsonify(list(cat_id_to_cat_obj.values()))


if __name__ == "__main__":
    init()
    app.run()
