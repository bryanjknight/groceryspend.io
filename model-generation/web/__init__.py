import json
import os

from flask import abort, Flask, request, jsonify

import joblib
import yaml

from training.categorize.train import CategorizeTuple

app = Flask(__name__)

cat_tuple: CategorizeTuple = joblib.load(os.environ['CAT_MODEL'])
cat_id_to_label = {}
with open(os.environ['CAT_YAML']) as f:
    y = yaml.safe_load(f)
    cat_id_to_label = { i['id']: i['name'] for i in y['categories']}


def init():
    pass

@app.route("/categorize", methods=["POST"])
def categorize():
    items = request.json
    if not items or len(items) == 0:
        return abort(400)

    item_features = cat_tuple.tfidf.transform(items)
    predictions = cat_tuple.model.predict(item_features)
    retval = { items[i]: cat_id_to_label[cat_tuple.id_to_cat[predictions[i]]] for i in range(len(predictions)) }
    return jsonify(retval)

@app.route("/categories", methods=["GET"])
def all_categories():
    return jsonify(cat_id_to_label)

if __name__ == "__main__":
    init()
    app.run()
