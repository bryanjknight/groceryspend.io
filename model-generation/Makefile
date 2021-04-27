
train:
  PYTHONPATH=. python training/run.py

run:
	pipenv run start

docker:
	docker build -t groceryspend.io/predict .

docker-run:
	docker run -it --rm -p 5000:5000 groceryspend.io/predict 

clean:
	rm -rf output/

local-train-and-run: | train run

docker-train-and-run: | train docker docker-run
