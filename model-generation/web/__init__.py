import json
import os

from flask import abort, Flask, request, jsonify, make_response, redirect, Response
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn import svm

import joblib
import numpy as np
import pandas as pd
import sklearn

from training.categorize.train import CategorizeTuple

app = Flask(__name__)

cat_tuple: CategorizeTuple = joblib.load(os.environ['CAT_MODEL'])

def init():
    pass

@app.route("/categorize", methods=["POST"])
def categorize():
    items = request.json
    if not items or len(items) == 0:
        return abort(400)

    item_features = cat_tuple.tfidf.transform(items)
    predictions = cat_tuple.model.predict(item_features)
    retval = { items[i]: cat_tuple.id_to_cat[predictions[i]] for i in range(len(predictions)) }
    return jsonify(retval)


if __name__ == "__main__":
    init()
    app.run()
