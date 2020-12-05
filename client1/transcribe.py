from __future__ import print_function
from urllib.request import urlopen
import json
import time
import boto3
from botocore.client import Config
import key

import requests

import datetime
import uuid
import os
import urllib.request

ACCESS_KEY_ID = key.KEY
ACCESS_SECRET_KEY = key.SECRET
BUCKET_NAME = 'tarun123456'

# BUCKET_NAME = os.environ['BUCKET_NAME']
s3 = boto3.resource('s3')


def transcription(url):
    # Download mp3 file to local server storage from url
    r = requests.get(url)

    with open("trial.mp3", 'wb') as f:
        f.write(r.content)
    a = "trial.mp3"

    # Upload the mp3 file to s3

    s3 = boto3.client(
        's3',
        aws_access_key_id=ACCESS_KEY_ID,
        aws_secret_access_key=ACCESS_SECRET_KEY,
        config=Config(signature_version='s3v4')
    )

    # Upload
    s3.upload_file(a, "tarun123456", "coen.mp3")
    os.remove(a)

    # Transcribing: setting up AWS parameters
    # New input - start
    transcribe = boto3.client('transcribe',
                              region_name='us-east-1',
                              aws_access_key_id=ACCESS_KEY_ID,
                              aws_secret_access_key=ACCESS_SECRET_KEY,
                              config=Config(signature_version='s3v4'))

    # New input - ends

    job_name = str(uuid.uuid4())
    job_uri = "https://s3.amazonaws.com/tarun123456/coen.mp3"  # https://s3.amazonaws.com/tarun123456/+"a"
    transcribe.start_transcription_job(
        TranscriptionJobName=job_name,
        Media={'MediaFileUri': job_uri},
        MediaFormat='mp3',
        LanguageCode='en-US'
    )
    while True:
        status = transcribe.get_transcription_job(TranscriptionJobName=job_name)
        if status['TranscriptionJob']['TranscriptionJobStatus'] in ['COMPLETED', 'FAILED']:
            uri = status['TranscriptionJob']['Transcript']['TranscriptFileUri']
            # print(uri)
            content = urllib.request.urlopen(uri).read().decode('UTF-8')

            # print(json.dumps(content))

            data = json.loads(content)
            transcribed_text = data['results']['transcripts'][0]['transcript']
            return transcribed_text
            # object = s3.Object('tarun123456',
            #     job_name+"_Output.txt",
            #     aws_access_key_id=ACCESS_KEY_ID,
            #     aws_secret_access_key=ACCESS_SECRET_KEY,
            #     config=Config(signature_version='s3v4'))
            # object.put(Body=transcribed_text)
            break
        print("Not ready yet...")
        time.sleep(5)


"""print(status)
with open('asrOutput.json','r') as dat:
    text =json.load(dat)
print(text)"""
