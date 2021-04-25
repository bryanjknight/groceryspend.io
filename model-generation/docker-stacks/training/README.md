1. `docker pull jupyter/scipy-notebook:latest`
1. Run with the following command in the root project directory: `docker run -it -p 8888:8888 -v $PWD/notebooks:/home/jovyan/work jupyter/scipy-notebook:latest`
1. Click the link to log into Jupyter
1. Files saved to `~/work` in the Jupyter notbook will be saved onto the local file system