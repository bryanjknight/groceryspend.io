# from https://towardsdatascience.com/multi-class-text-classification-with-scikit-learn-12f1e60e0a9f

import collections
import glob
import logging

import pandas as pd

from sklearn import metrics
from sklearn.svm import LinearSVC
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.model_selection import train_test_split

CategorizeTuple = collections.namedtuple(
    "CategorizeTuple", ["id_to_cat", "tfidf", "model"]
)

LOG = logging.getLogger(__name__)


def _read_csvs_into_df(data_dir: str) -> pd.DataFrame:
    """Assumes that all csvs have same format and columns"""

    all_files = glob.glob(data_dir + "/*.csv")

    df = pd.concat((pd.read_csv(f) for f in all_files))
    return df


def train(data_dir: str, min_accuracy: float) -> CategorizeTuple:
    """Using a directory of csvs, train a model to categorize an item
    to a grocery category"""

    LOG.info("Reading csvs into dataframe")
    df = _read_csvs_into_df(data_dir)

    # create mappings between category and category id
    LOG.info("Creating id to category mapping")
    df["category_id"] = df["Category"].factorize()[0]
    category_id_df = (
        df[["Category", "category_id"]].drop_duplicates().sort_values("category_id")
    )
    id_to_category = dict(category_id_df[["category_id", "Category"]].values)

    # create tfdif vectorizer
    LOG.info("Creating TFIDF")
    tfidf = TfidfVectorizer(
        sublinear_tf=True,
        min_df=5,
        norm="l2",
        encoding="latin-1",
        ngram_range=(1, 6),
        stop_words="english",
    )
    features = tfidf.fit_transform(df.Item).toarray()
    labels = df.category_id

    # split data int train and test sets
    # through testing, we determined LinearSVC was the best performing
    # TODO: auto-evaluate different models to verify best performing model
    LOG.info("Fitting model")
    model = LinearSVC()
    x_train, x_test, y_train, y_test, _, _ = train_test_split(
        features, labels, df.index, test_size=0.33, random_state=0
    )
    model.fit(x_train, y_train)

    # test to verify it meets our criteria
    # TODO: should this be separated? Seems easier since we have the
    #       test features already here
    LOG.info("Running model tests")
    y_pred = model.predict(x_test)
    report = metrics.classification_report(
        y_test, y_pred, target_names=df["Category"].unique(), output_dict=True
    )

    actual_accuracy = report["accuracy"]
    LOG.info(f"Observed accuracy: {actual_accuracy}")
    if actual_accuracy < min_accuracy:
        raise Exception(
            f"Training resulted in a model with acc of {actual_accuracy}, min is {min_accuracy}"
        )

    return CategorizeTuple(id_to_category, tfidf, model)
