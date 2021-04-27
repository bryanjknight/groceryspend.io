"""Module for training a categorize model"""
from typing import Any
import logging

import joblib

from training.categorize.train import CategorizeTuple, train
from training.categorize.fetch import stopandshop
# from training.categorize.fetch.hayd.parser import Parser

LOG = logging.getLogger(__name__)

def fetch_all_data(temp_dir: str):
    """ fetch all data for training categorize model """
    LOG.info("Fetching Stop 'N Shop Data")
    stopandshop.run(temp_dir)

    # TODO: run HAYD export as well


def persist_objects(obj: Any, model_dir: str):
    LOG.info("Persisting categorize tuple")
    joblib.dump(obj, f"{model_dir}/categorize_tuple.pkl")


def clean_up(data_dir: str):
    pass


def run(data_dir: str, model_dir: str, min_accuracy: float):
    # fetch all data
    fetch_all_data(data_dir)

    # train model
    LOG.info("Training categorize model")
    cat_tuple: CategorizeTuple = train(data_dir, min_accuracy)

    # persist objects to filesystem
    persist_objects(cat_tuple, model_dir)

    # clean up
    clean_up(data_dir)
