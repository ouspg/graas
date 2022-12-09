from flask import Flask, render_template, request, make_response, g
import os
import socket
import random
import json
import logging

hostname = socket.gethostname()

app = Flask(__name__)

gunicorn_error_logger = logging.getLogger('gunicorn.error')
app.logger.handlers.extend(gunicorn_error_logger.handlers)
app.logger.setLevel(logging.INFO)

MESSAGE1 = "Hi Bob! Yes! I am well. I was just watching this lol https://www.youtube.com/watch?v=AQDCe585Lnc"
MESSAGE2 = "https://i.imgur.com/04HJwHg.gif"
MESSAGE3 = "My lips are sealed. What a good news. Our stock is going through the roof!"

@app.route("/", methods=['POST','GET'])
def receive_message():

    if request.method == 'POST':
        message = request.json.get("message")
        resp_message = ""
        if message.startswith("Hey Alice"):
            resp_message = MESSAGE1
        elif message.startswith("luckily this"):
            resp_message = MESSAGE2
        elif message.startswith("Yes. It is time."):
            resp_message = MESSAGE3
        elif message in ["Let me check that...", "lol we famous..", "Yes."]:
            return make_response("", 200)
        else:
            return make_response("What u talking about", 400)
        resp = make_response(
            json.dumps(
                        {"message": resp_message, "confidentality": "TOP SECRET"}
                    ),
                    200,
        )
        resp.headers['User-Agent'] = "Alice's Messenger"
        return resp
    else:
        return make_response("Hi, I am Alice.", 200)


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=8000, debug=True, threaded=True)

# https://imgs.xkcd.com/comics/alice_and_bob.png