Python ML
===

Given that ML tools are mostly written in Python, this is a spike to quickly assess how to categorize all of our products into usable categories

Install
---
* `pyenv` to install python 3.9
* Set up your terminal to use the .pyenv directory 
* `pyenv install 3.9.4` (change minor version accordingly)
* `pipenv install --pre`
* `PIPENV_VENV_IN_PROJECT="enabled" pipenv install ` to create the pip environment


Train
---
1. `PYTHONPATH=. python training/run.py`

Deploy as Flask API
---
1. `FLASK_APP=web python -m flask run`
