Python ML
===

Given that ML tools are mostly written in Python, this is a spike to quickly assess how to categorize all of our products into usable categories

Install
---
* `pyenv` to install python 3.9
* Set up your terminal to use the .pyenv directory 
* `pyenv install 3.9.4` (change minor version accordingly)
* `PIPENV_VENV_IN_PROJECT="enabled" pipenv install ` to create the pip environment


Notes
---
* Using https://scikit-learn.org/stable/tutorial/text_analytics/working_with_text_data.html

* I need training data
  * https://helpatyourdoor.org/wp-content/uploads/2018/02/2018-Store-Brand-Catalog.pdf
* I need to test the trained model with other data



TensorFlow Approach
===

* `docker pull tensorflow/tensorflow:latest-jupyter`
  * get the latest tensorflow with jupyter support

* Build a new image with the additional libraries we need


* `docker run -it -p 8888:8888 -v $PWD/tf:/ bk-test`
  * start an instance, save notebooks to `./notebooks`