""" Run training for all models """
import logging
import os
import sys

import training.categorize.run as cat_run

OUTPUT_FORMAT = "output/{model}/{bucket}"


def _setup_output_dirs(model: str):
    tmp_data_dir = OUTPUT_FORMAT.format(model=model, bucket="data")
    tmp_model_dir = OUTPUT_FORMAT.format(model=model, bucket="model")
    os.makedirs(tmp_data_dir, exist_ok=True)
    os.makedirs(tmp_model_dir, exist_ok=True)
    return (tmp_data_dir, tmp_model_dir)


if __name__ == "__main__":

    # setup logger
    logging.basicConfig(stream=sys.stdout, level=logging.INFO)

    # run training for categorize
    print("Training categorize model")
    data_dir, model_dir = _setup_output_dirs("categorize")
    cat_run.run(data_dir, model_dir, 0.75)
