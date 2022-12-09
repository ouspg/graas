import requests
import time
from subprocess import PIPE, run

MESSAGE1 = "Hey Alice! It's been busy times. How it is going?"
MESSAGE2 = "Let me check that..."
MESSAGE3 = "lol we famous.."
MESSAGE4 = "luckily this conversation is safe. I have a business suggestion."
result = run("/usr/share/secret", stdout=PIPE, stderr=PIPE, universal_newlines=True)
SECRET = result.stdout.strip()
MESSAGE5 = f"Yes. It is time. Time to make our company public. I found a major investor. He is {SECRET}. Don't tell anyone, otherwise we might lose him at this point."
MESSAGE6 = "Yes."

url = "http://alice"

time.sleep(15)

with requests.Session() as s:
    s.headers.update({'User-Agent': 'Bob\'s Messenger'})
    message_base = {"message": "", "confidentality": "TOP SECRET"}    
    message_base["message"] = MESSAGE1
    resp = s.post(url, json=message_base)
    while True:
        resp_message = resp.json().get("message")
        if resp_message:
            if resp_message.startswith("Hi Bob!"):
                time.sleep(3)
                message_base["message"] = MESSAGE2
                resp = s.post(url, json=message_base)
                if resp.status_code == 200:
                    s.get("https://www.youtube.com/watch?v=AQDCe585Lnc")
                    time.sleep(5)
                message_base["message"] = MESSAGE3
                resp = s.post(url, json=message_base)
                time.sleep(5)
                message_base["message"] = MESSAGE4
                resp = s.post(url, json=message_base)

            if resp_message.startswith("https"):
                s.get("https://i.imgur.com/04HJwHg.gif")
                time.sleep(3)
                message_base["message"] = MESSAGE5
                resp = s.post(url, json=message_base)
            if resp_message.startswith("My lips are sealed."):
                time.sleep(3)
                message_base["message"] = MESSAGE6
                resp = s.post(url, json=message_base)
                time.sleep(2)
                break
        