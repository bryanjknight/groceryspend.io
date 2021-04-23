1. Build from this directory: `docker build -t groceryspend.io/training .`. Ideally this is one time, but may require rebuilding as new packages are needed
1. Run with the following command in the root project directory: `docker run -it -p 8888:8888 -v $PWD/notebooks:/tf/notebooks groceryspend.io/training`
1. Click the link to log into Jupyter