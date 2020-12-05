# Color palettes: #061323 #798478 #FFFEFC
# Set Flask to debug mode :$env:FLASK_ENV = "development"
from flask import Flask, render_template, request, url_for, redirect, session, abort
from response_processor import keyword_processor
import requests
import json
import time
import os

import checkemail
# import transcribe
import topics
import recomender

application = app = Flask(__name__)
app.secret_key = "coen6313"
# main_url = "http://podcast.ca-central-1.elasticbeanstalk.com"  # "http://localhost:6313"


main_url = "http://localhost:5001"

#
# @app.route("/")
# def search():
#     return render_template("search.html", username=session.get('username'), login_token=session.get('token'))

@app.route("/")
def search():
    return render_template("search.html", username=session.get('username'), login_token=session.get('token'))


@app.route("/result", methods=["GET", "POST"])
def result():
    global main_url
    name = request.cookies.get('token')
    print(name)
    keyword = str(request.form.get("searchkey"))
    newkeyword = keyword_processor(keyword)
    url = 'https://listen-api.listennotes.com/api/v2/search?q=' + newkeyword + '&sort_by_date=0&type=episode&offset=0' \
                                                                               '&len_min=3&len_max=2&genre_ids=68' \
                                                                               '%2C82&published_before=1580172454000' \
                                                                               '&published_after=0&only_in=title' \
                                                                               '%2Cdescription&language=English' \
                                                                               '&safe_mode=0 '
    headers = {
        'X-ListenAPI-Key': 'ce8e7ce414414764be7159d0aeecdb16',
    }

    # comment out the sending request to style and modify result page
    response = requests.request('GET', url, headers=headers).json()

    # with open("test.json", "r") as json_file:
    #     response = json.load(json_file)

    if response["results"] == []:
        return abort(404)
    else:
        # response=json.dumps(response)
        # results = response_processor(response)
        return render_template("result.html", keyword=keyword, response=response, username=session.get('username'),
                               login_token=session.get('token'))


@app.route("/sign_in", methods=["POST", "GET"])
def sign_in():
    global main_url
    if request.method == "POST":
        url = main_url + "/login"

        # Processing if the input is an email or username:
        email_username = str(request.form.get("email_username"))

        if checkemail.check(email_username):
            params = {"Email": str(request.form.get("email_username")),
                      "Password": str(request.form.get("password"))}
        else:
            params = {"Username": str(request.form.get("email_username")),
                      "Password": str(request.form.get("password"))}

        r = requests.request("POST", url, data=params).json()

        if r["success"]:
            session['token'] = r['token']
            # session['email'] = r['user']['email']
            # session['username'] = r['user']['username']
            # hed = {'Authorization': 'Bearer ' + session.get('token')}
            # r = requests.get("http://localhost:6313/login",headers=hed).json()
            # session['fullname'] = r['full_name']
            # session['favorite'] = r['favorite_podcast bson:']
            return redirect(url_for('account'))
        else:
            return r["message"]
        # If log in succeeds -> global login_token = return
        # Else returns log_in page with error message (must contain jinja if to process the message)

    return render_template("signin.html")


@app.route("/sign_up", methods=["POST", "GET"])
def sign_up():
    global main_url
    if request.method == "POST":
        if str(request.form.get("password")) != str(request.form.get("confirm_password")):
            return render_template("signup.html", passwordmatch=True)
        else:
            url = main_url + "/signup"
            body = {"Username": str(request.form.get("username")),
                    "FullName": str(request.form.get("fullname")),
                    "Email": str(request.form.get("email")),
                    "Password": str(request.form.get("password")),
                    "ConfirmPassword": str(request.form.get("confirm_password"))}
            r = requests.request("POST", url, data=body).json()
            if r["success"]:
                return render_template("signin.html")
            else:
                return r["message"]

    return render_template("signup.html", passwordmatch=False)


@app.route("/sign_out")
def sign_out():
    global main_url
    session.pop('token', None)
    session.pop('email', None)
    session.pop('username', None)
    session.pop('fullname', None)
    session.pop('favorite', None)
    return redirect(url_for("search"))


@app.route("/favorite", methods=["POST"])
def favorite():
    global main_url
    if 'token' not in session:
        return redirect(url_for('sign_in'))
    else:
        url = main_url + "/favorite"
        hed = {'Authorization': 'Bearer ' + session.get('token')}
        body = {"ID": request.form["id"],
                "Name": request.form["name"],
                "ImageUrl": request.form["image"]}
        r = requests.request("POST", url, data=body, headers=hed)
        return redirect(url_for("account"))


@app.route("/podcast_details", methods=['POST'])
def podcast_details():
    global main_url
    podcast = {}
    if request.form.get('audio') is not None:
        # accessing podcast details from result page
        podcast['id'] = request.form.get('id')
        podcast['audio'] = request.form.get('audio')
        podcast['episode'] = request.form.get('episode_title')
        podcast['podcast'] = request.form.get('podcast_title')
        podcast['image'] = request.form.get('image')
        # Finding the recommended podcasts
        recommender_input_id = request.form.get('podcast_id')
        headers = {
            'X-ListenAPI-Key': 'ce8e7ce414414764be7159d0aeecdb16',
        }
        url = 'https://listen-api.listennotes.com/api/v2/podcasts/' + recommender_input_id + '?next_episode_pub_date=1479154463000&sort=recent_first'
        features = requests.request('GET', url, headers=headers).json()
        rec_podcast_ids = recomender.recommender(features)
        rec_podcast_array = []
        for i in rec_podcast_ids:
            url = 'https://listen-api.listennotes.com/api/v2/podcasts/' + i + '?next_episode_pub_date=1479154463000&sort=recent_first'
            rec_podcast_array.append(requests.request('GET', url, headers=headers).json())
    else:
        # accessing podcast details from account page
        podcast['id'] = request.form.get('podcast_id')
        url = 'https://listen-api.listennotes.com/api/v2/episodes/' + podcast['id']
        headers = {
            'X-ListenAPI-Key': 'ce8e7ce414414764be7159d0aeecdb16',
        }
        response = requests.request('GET', url, headers=headers).json()
        podcast['audio'] = response["audio"]
        podcast['episode'] = response["title"]
        podcast['podcast'] = response["podcast"]["title"]
        podcast['image'] = response["image"]
        features = response["podcast"]
        rec_podcast_ids = recomender.recommender(features)
        rec_podcast_array = []
        for i in rec_podcast_ids:
            url = 'https://listen-api.listennotes.com/api/v2/podcasts/' + i + '?next_episode_pub_date=1479154463000&sort=recent_first'
            rec_podcast_array.append(requests.request('GET', url, headers=headers).json())
    if 'token' not in session:
        return redirect(url_for('sign_in'))
    else:
        # #Testing module
        # podcast['id'] = "8895bb5053d948c3bedf37d80394f050"

        # print(session.get('token'))
        # print("The suppose ID is: "+podcast['id'])
        url = main_url + "/gettopics"
        hed = {'Authorization': 'Bearer ' + session.get('token')}
        params = {"id": podcast['id']}
        r = requests.get(url, params=params, headers=hed).json()
        # print(r)
        audio = podcast['audio']
        if r['success'] == False:
            # #Testing module
            # with open('podcastdata.json') as f:
            #     data= json.load(f)
            print("we are getting topics")
            # timestamp needs to change to comprehend code
            print(audio)
            speech_to_text_data, topics_data = topics.getting_topics(audio)

            timestamp = topics_data
            transcript = speech_to_text_data['transcript']
            # #Test data:
            # timestamp = data['time_stamp']
            # transcript = data['transcript']
            # send the new podcast entry to the backend
            url2 = main_url + "/savetopics"

            hed = {'Authorization': 'Bearer ' + session.get('token'), 'Content-type': 'application/json',
                   'Accept': 'text/plain'}
            body = {
                "name": podcast['episode'],
                "id": podcast['id'],
                "image_url": podcast['image'],
                "transcript": transcript,
                "time_stamp": timestamp
            }
            # print(hed)
            # with open('podcastdata.json', 'w') as f:
            #     json.dump(body, f)
            # print(json.load(f))
            r2 = requests.request("POST", url2,
                                  data=json.dumps(body), headers=hed).json()
            # print(r2)

        else:
            # Testing module
            # with open('podcastdata.json') as f:
            #     data2 = json.load(f)
            # print(r)
            data = r['data']
            timestamp = data['topics']
            # timestamp = data2['time_stamp']
            # timestamp = data2['topics']
            transcript = data['transcript']

        return render_template("podcast_details.html", login_token=session.get('username'), transcript=transcript,
                               podcast=podcast, timestamp=timestamp, rec_podcast_array=rec_podcast_array)


@app.route("/guser")
def guser():
    global main_url
    user_id = request.args.get('user_id', type=str)
    url = main_url + "/guser?user_id=" + user_id
    r = requests.get(url).json()
    if r["success"]:
        session['token'] = r['token']
        return redirect(url_for('account'))
    else:
        abort(404)


@app.route("/account")
def account():
    global main_url
    # Processing oauth log in
    details = request.cookies.get('user_google')
    if details:
        detail = details.split("|")
        email = detail[0]
        token = detail[1]
        session['token'] = token

    hed = {'Authorization': 'Bearer ' + session.get('token')}
    r = requests.get(main_url + "/getuser", headers=hed).json()
    session['email'] = r['email']
    session['username'] = r['username']
    session['fullname'] = r['full_name']
    session['favorite'] = r['favorite_podcast bson:']

    if 'token' not in session:
        return redirect(url_for('sign_in'))
    else:

        # Should have any user information here for sending requests

        # Replace with get request to back-end to get the info
        user = {'username': session.get('username'), 'fullname': session.get('fullname'), 'email': session.get('email')}

        # Replace with get request to back-end to get the info
        favPodcast = session.get('favorite')
        # print(session.get('token'))
        return render_template("account.html", username=session.get('username'), login_token=session.get('token'),
                               favPodcast=favPodcast, user=user)


if __name__ == "__main__":
    app.run()
