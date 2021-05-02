Python ML
===

Given that ML tools are mostly written in Python, this is a spike to quickly assess how to categorize all of our products into usable categories

Install
---
* `pyenv` to install python 3.7
* Set up your terminal to use the .pyenv directory 
* `pyenv install 3.7.10` (change minor version accordingly)
* The `.python-version` file will tell `pyenv` to use `3.7.10`
* `pip install pipenv`
* `PIPENV_VENV_IN_PROJECT="enabled" pipenv install ` to create the pip environment


Train
---
1. `make train`

Deploy as Flask API
---
1. `FLASK_APP=web python -m flask run`
