FROM python:3.9
WORKDIR /app
COPY . /app
RUN apt-get -y update
RUN pip3 install -r requirements.txt
RUN python -m nltk.downloader all
#Expose the required port
EXPOSE 5000
#Run the command
ENTRYPOINT [ "python3" ]
CMD [ "app.py" ]