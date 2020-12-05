import os, pathlib2, requests, io
from sys import platform
from heapq import nlargest
from mutagen.mp3 import MP3
from google.cloud import speech_v1p1beta1 as speech
from google.oauth2 import service_account
from google.cloud import language_v1
from google.cloud import storage

import uuid

import nltk
import time
from urllib.request import Request, urlopen




def download_uploa_GCS(url, keypath):
    # Create a unique name for audio file
    audio_name = str(uuid.uuid4()) + '.mp3'
    req = Request(url, headers={'User-Agent': 'Mozilla/5.0'})

    web_byte = urlopen(req).read()

    # filename = wget.download(url)
    # # Set up credentials from local keypath
    credentials = service_account.Credentials.from_service_account_file(keypath)
    # Initiate google cloud storage service
    storage_client = storage.Client(credentials=credentials)
    # Provide bucket name to upload the audio
    bucket = storage_client.bucket('coen6313project')
    # Provide audio file name for its presence on the cloud (it is the same as local file)
    blob = bucket.blob(audio_name)
    # Upload the file to the cloud
    blob.upload_from_string(web_byte, content_type='audio/mpeg')
    # Return the URI of the file of the cloud
    gcs_URI = 'gs://coen6313project/' + audio_name

    return gcs_URI

def download_upload_GCS(url, keypath):
    # Create a unique name for audio file
    audio_name = str(uuid.uuid4()) + '.mp3'
    # Get the audio file content from the url
    r = requests.get(url)
    # Write the audio file locally to client's server
    with open(audio_name, 'wb') as f:
        f.write(r.content)

    # Get the sample rate of the downloaded audio
    audio_info = MP3(audio_name).info
    sample_rate = audio_info.sample_rate

    # Set up a file path to upload the file to google cloud
    path = pathlib2.Path(__file__).parent.absolute()
    filepath = audio_name
    # Set up credentials from local keypath
    credentials = service_account.Credentials.from_service_account_file(keypath)
    # Initiate google cloud storage service
    storage_client = storage.Client(credentials=credentials)
    # Provide bucket name to upload the audio
    bucket = storage_client.bucket('coen6313project')
    # Provide audio file name for its presence on the cloud (it is the same as local file)
    blob = bucket.blob(audio_name)
    # Upload the file to the cloud
    blob.upload_from_filename(filepath)
    # Return the URI of the file of the cloud
    gcs_URI = 'gs://coen6313project/' + audio_name
    # Remove audiofile from local storage:
    os.remove(audio_name)
    return gcs_URI, sample_rate
# dbe1d40543d34d92aa2e1bc32883c6

# Google speech to text part-------------------------------------
# The output of this function should be a list of dictionary that contains all the start time and all the words
# Then the output is saved as a json file to be further manipulated down the road.
def speech_to_text(gcs_URI, keypath):
    # Reference: https://cloud.google.com/speech-to-text/docs/async-recognize
    # Set up credentials from local keypath
    G = 'https://www.listennotes.com/e/p/ea09b575d07341599d8d5b71f205517b/'
    credentials = service_account.Credentials.from_service_account_file(keypath)
    audio = speech.RecognitionAudio(uri=gcs_URI)
    config = speech.RecognitionConfig(
        language_code="en-US",
        enable_automatic_punctuation=True,
        enable_word_time_offsets=True,
        encoding=speech.RecognitionConfig.AudioEncoding.MP3,
        sample_rate_hertz=16000,
    )

    client = speech.SpeechClient(credentials=credentials)
    operation = client.long_running_recognize(config=config, audio=audio)
    print("Waiting for operation to complete...")
    response = operation.result()
    i = 1
    sentence = ''
    transcript_all = ''
    start_time_offset = []
    # Building a python dict (contains start time and words) from the response:
    for result in response.results:
        best_alternative = result.alternatives[0]
        transcript = best_alternative.transcript
        if i == 1:
            transcript_all = transcript
        else:
            transcript_all += " " + transcript
        i += 1
        # Getting timestamps
        for word in best_alternative.words:
            start_s = word.start_time.total_seconds()
            word = word.word
            if sentence == '':
                sentence = word
                sentence_start_time = start_s
            else:
                sentence += ' ' + word
                if '.' in word:
                    start_time_offset.append({'time': sentence_start_time, 'sentence': sentence})
                    sentence = ''
    speech_to_text_data = {'transcript': transcript_all, 'timestamps': start_time_offset}
    print('Finish transcription.')
    return speech_to_text_data


# Getting key phrases for texts:

def top_sentence(text):
    # Reference: https://medium.com/analytics-vidhya/simple-text-summarization-using-nltk-eedc36ebaaf8
    # This is biased to long sentences

    from string import punctuation
    word_tokenize = nltk.tokenize.word_tokenize
    sent_tokenize = nltk.tokenize.sent_tokenize
    stopwords = nltk.corpus.stopwords
    tokens = word_tokenize(text)
    stop_words = stopwords.words('english')
    punctuation = punctuation + '\n'

    word_frequencies = {}
    for word in tokens:
        if word.lower() not in stop_words:
            if word.lower() not in punctuation:
                if word not in word_frequencies.keys():
                    word_frequencies[word] = 1
                else:
                    word_frequencies[word] += 1

    max_frequency = max(word_frequencies.values())

    for word in word_frequencies.keys():
        word_frequencies[word] = word_frequencies[word] / max_frequency

    sent_token = sent_tokenize(text)

    sentence_scores = {}
    for sent in sent_token:
        sentence = sent.split(" ")
        for word in sentence:
            if word.lower() in word_frequencies.keys():
                if sent not in sentence_scores.keys():
                    sentence_scores[sent] = word_frequencies[word.lower()]
                else:
                    sentence_scores[sent] += word_frequencies[word.lower()]

    # Final summary for each part will contain 2 sentences with heaviest weight
    select_length = 1  # int(len(sent_token)*0.05)

    summary = nlargest(select_length, sentence_scores, key=sentence_scores.get)

    final_summary = [word for word in summary]
    summary = ' '.join(final_summary)
    return summary


# Partitioning the transcript using NLTK package:
def partitioning_transcript(podcast_dictionary):
    # Reference: https://www.nltk.org/_modules/nltk/tokenize/texttiling.html
    print('Starting partitioning')
    text = podcast_dictionary['transcript']
    timestamps = podcast_dictionary['timestamps']
    topics = []
    text_prep = text.replace('.', '.\n\n\n')
    tt = nltk.tokenize.texttiling.TextTilingTokenizer(w=5, k=10)

    segmented_text = tt.tokenize(text_prep)
    new_text = ''
    for item in segmented_text:
        item = item.replace('\n\n\n', '')
        item += ' \n\n\n'
        new_text += item
    tt = nltk.tokenize.texttiling.TextTilingTokenizer(w=50, k=2)
    segmented_text = tt.tokenize(new_text)

    for item in segmented_text:
        item = item.replace('\n\n\n', '')
        timestamps_array = []
        # Transferring timestamps to each partition:
        for timestamp in timestamps:
            if timestamp['sentence'] in item:
                timestamps_array.append(timestamp['time'])
        sec = int(min(timestamps_array))
        ty_res = time.gmtime(sec)
        if sec >= 3600:
            res = time.strftime("%H:%M:%S", ty_res)
        else:
            res = time.strftime("%M:%S", ty_res)
        topics.append({"time_stamp": res, "topic": top_sentence(item)})
    print('Finishing partitioning')
    return topics


def getting_topics(url):
    # $env:GOOGLE_APPLICATION_CREDENTIALS="C:\Users\OS\COEN6313_Assignment2\ivory-hallway-296414-c2242781fe7f.json"

    path = pathlib2.Path(__file__).parent.absolute()
    print(path)
    print(os.getcwd())
    keypath = 'ivory-hallway-296414-7879422a328c.json'
    if os.getenv('GOOGLE_APPLICATION_CREDENTIALS') is None:
        os.environ['GOOGLE_APPLICATION_CREDENTIALS'] = keypath

    gcs_URI = download_uploa_GCS(url, keypath)

    speech_to_text_data = speech_to_text(gcs_URI, keypath)

    # out_file = open("transcript2.json", "w")

    # json.dump(speech_to_text_data, out_file, indent = 6)

    # out_file.close()
    # in_file = open("transcript2.json",)
    # speech_to_text_data = json.load(in_file)
    topics = partitioning_transcript(speech_to_text_data)

    # out_file = open("topics.json", "w")

    # json.dump(topics, out_file, indent = 6)

    # out_file.close()
    return speech_to_text_data, topics
# url = 'https://www.listennotes.com/e/p/5ee68888b5484cabb99c868b26e81b16/'
# speech_to_text_data, topics_data=getting_topics(url)
# print(speech_to_text_data)
# print('\n')
# print(topics_data)

# Google analyzing content and partitioning part--------------------------------------
# #Rules:
# # 1. partitions must be at least 3 minutes long
# # 2. partitions can be up to 10-20 minutes long
# # 3. partitions are separated by the density of the entities from google natural language API:
# # 3*: The entities can density can have a certain distribution wrt time, and it could be represented by a discrete
# # distribution function. The cut point is the at the middle between the the valley after the 3 minutes mark and the peak
# # after it where the sentences are different.
# #Analyzing content starts:
# def sample_analyze_entities(text_content):
#     """
#     Analyzing Entities in a String

#     Args:
#       text_content The text content to analyze
#     """

#     client = language_v1.LanguageServiceClient(credentials=credentials)

#     # text_content = 'California is a state.'

#     # Available types: PLAIN_TEXT, HTML
#     type_ = language_v1.Document.Type.PLAIN_TEXT

#     # Optional. If not specified, the language is automatically detected.
#     # For list of supported languages:
#     # https://cloud.google.com/natural-language/docs/languages
#     language = "en"
#     document = {"content": text_content, "type_": type_, "language": language}

#     # Available values: NONE, UTF8, UTF16, UTF32
#     encoding_type = language_v1.EncodingType.UTF8

#     response = client.analyze_entities(request = {'document': document, 'encoding_type': encoding_type})

#     # Loop through entitites returned from the API
#     for entity in response.entities:
#         global data
#         # Getting the entity name
#         name = entity.name
#         # Get entity type, e.g. PERSON, LOCATION, ADDRESS, NUMBER, et al
#         entity_type = language_v1.Entity.Type(entity.type_).name
#         # Get the salience score associated with the entity in the [0, 1.0] range
#         score = entity.salience
#         # Loop over the metadata associated with entity. For many known entities,
#         # the metadata is a Wikipedia URL (wikipedia_url) and Knowledge Graph MID (mid).
#         # Some entity types may have additional metadata, e.g. ADDRESS entities
#         # may have metadata for the address street_name, postal_code, et al.
#         # for metadata_name, metadata_value in entity.metadata.items():
#         #     print(u"{}: {}".format(metadata_name, metadata_value))

#         # Loop over the mentions of this entity in the input document.
#         # The API currently supports proper noun mentions.
#         for mention in entity.mentions:
#             mention_text = mention.text.content
#             # Add the a field entity of type boolean inside to timestamps dict that tells if a word is an entity or not
#             # First check if mention text is a group of word or not
#             if ' ' in mention_text:
#                 single_text = mention_text.split(' ')
#             else:
#                 single_text = [mention_text]
#             print(single_text)
#             # Then for each word in the group of word
#             for text in single_text:
#             # Find the dict in timestamp array which has the value equals to mention_text. The value corresponds to
#             # the key "word" of the dict.
#             # something like :
#             # test_list = [{'gfg' : 2, 'is' : 4, 'best' : 6},
#             # {'it' : 5, 'is' : 7, 'best' : 8},
#             # {'CS' : 10}]
#             # res = next((sub for sub in test_list if sub['is'] == 7), None)
#                 res = next((dict_element for dict_element in data['timestamps'] if
#                     dict_element['word'].strip('.,?').replace('\'s','') == str(text).strip('.,?').replace('\'s','')
#                     or (str(entity_type) == "NUMBER" and str(text)+'s' == dict_element['word'])
#                     ), None)
#             # AND Change the value of "entity" key to true:
#                 if res is not None:
#                     data['timestamps'][data['timestamps'].index(res)]['entity'] = True

#             # Get the mention type, e.g. PROPER for proper noun
#             mention_type = language_v1.EntityMention.Type(mention.type_).name
#         global entity_list
#         entity_list.append(
#             {'name' : name,
#             'entity_type': entity_type,
#             'score': score,
#             'mention_text': mention_text,
#             'mention_type': mention_type})

#     # Get the language of the text, which will be the same as
#     # the language specified in the request or, if not specified,
#     # the automatically-detected language.

# in_file = open("transcript.json", )

# data = json.load(in_file)

# in_file.close()

# transcript = data['transcript']

# entity_list = []

# sample_analyze_entities(transcript)

# out_file = open("entity_list.json", "w")

# json.dump(entity_list, out_file, indent = 6)

# out_file.close()

# out_file = open("transcript.json", "w")

# json.dump(data, out_file, indent = 6)
