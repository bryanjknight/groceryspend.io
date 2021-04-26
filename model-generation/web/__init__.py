import json
import pickle

from flask import abort, Flask, request, jsonify, make_response, redirect, Response
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn import svm

import joblib
import numpy as np
import pandas as pd
import sklearn


app = Flask(__name__)

# TODO: make this more configurable

id_to_category = joblib.load("/Users/bknight/dev/groceryspend.io/monorepo/model-generation/notebooks/id_to_category.pkl")
categorize_model = joblib.load("/Users/bknight/dev/groceryspend.io/monorepo/model-generation/notebooks/model.pkl")
tfidf = joblib.load("/Users/bknight/dev/groceryspend.io/monorepo/model-generation/notebooks/tfidf.pkl")


def init():
    pass



@app.route("/categorize", methods=["POST"])
def categorize():
    items = request.json
    if not items or len(items) == 0:
        return abort(400)

    item_features = tfidf.transform(items)
    predictions = categorize_model.predict(item_features)
    retval = { items[i]: id_to_category[predictions[i]] for i in range(len(predictions)) }
    return jsonify(retval)


if __name__ == "__main__":
    init()
    app.run()
