from json import dumps
from sklearn.preprocessing import LabelEncoder
from sklearn.neighbors import NearestNeighbors
import pandas as pd
import numpy as np
import requests
import pickle
from six.moves.urllib.parse import urlencode, quote

ref_data = pd.read_excel("Book3.xlsx")

# cleaning and processing the data
ref_data["explicit_conten"] = ref_data["explicit_conten"].astype("int64")
ref_data.dropna(inplace=True)

print(ref_data)

# preparing data for training
tr_data = pd.DataFrame(
    {"episodes": ref_data["total_episodes"], "exp": ref_data["explicit_conten"], "gen": ref_data["genre_ids"]})

# creating instance of the recomender algorithm
nb = NearestNeighbors()

# training the algorithm
nb.fit(tr_data)

pkl_filename = "recommender.pkl"

with open(pkl_filename, 'wb') as file:
    pickle.dump(nb, file)
