from json import dumps
from sklearn.preprocessing import LabelEncoder
from sklearn.neighbors import NearestNeighbors
import pandas as pd
import numpy as np
import requests
import pickle
from six.moves.urllib.parse import urlencode, quote


def recommender(features):
    # Input - features is of podcast - Implemented at the podcast details page
    # features is a dict with keys : total_episodes, explicit_content, genre_ids

    # print(tr_data)

    # # Processing the user input
    # pod_input = (input("Search for any podcast :").lower()).strip()
    # pod_input = pod_input.replace(" ","%20")
    # print(pod_input)

    # # user searching REST api process
    # url = f'https://listen-api.listennotes.com/api/v2/search?q={pod_input}&sort_by_date=0&type=episode&offset=0&len_min=10&len_max=30&genre_ids=68%2C82&published_before=1580172454000&published_after=0&only_in=title%2Cdescription&language=English&safe_mode=0'
    # headers = {
    #         'X-ListenAPI-Key': 'ce8e7ce414414764be7159d0aeecdb16'}

    # # getting podcast id form user search
    # pre_response = requests.request('GET', url, headers=headers)
    # pre_response = pre_response.json()
    # pre_response_id = (pre_response["results"][0])["podcast"]["id"]

    # #print("pre_response",pre_response_id)

    # #getting features form user search podcast id_url = f'https://listen-api.listennotes.com/api/v2/podcasts/{
    # pre_response_id}?next_episode_pub_date=1479154463000&sort=recent_first'

    # features = requests.request('GET', id_url, headers=headers)
    # features = features.json()

    # #print("features",features)

    # This line is unknown
    ref_data = pd.read_excel("Training_data.xlsx")

    # cleaning and processing the data
    ref_data["explicit_conten"] = ref_data["explicit_conten"].astype("int64")
    ref_data.dropna(inplace=True)

    # Loading pretrained data
    pkl_filename = "recommender.pkl"

    with open(pkl_filename, 'rb') as file:
        nb = pickle.load(file)

    # processing the features of podcast and fitting them to algorithm
    if features:
        us1 = int(features["total_episodes"])
        us2 = int(features["explicit_content"])
        us3 = int(features["genre_ids"][-1])
        # us4 = int(response["listen_score"])
        ans = nb.kneighbors(np.array([[us1, us2, us3]]), return_distance=False)
        re1 = (ref_data.iloc[ans[0][0]])["ID"]  # change title to id
        re2 = (ref_data.iloc[ans[0][1]])["ID"]
        re3 = (ref_data.iloc[ans[0][2]])["ID"]

    # print(re1,re2,re3,"after",sep="\n")
    return [re1, re2, re3]

# def recommendation_system(pod_input):
#   Your code here
#   return epid1, epid2, epid3
