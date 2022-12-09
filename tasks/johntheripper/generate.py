import crypt
import os
import sys
import random

# Task description: Brute force the shadow file with John The Ripper


def password_from_seed():

    ###########

    SEED = os.environ.get('STUDENT_SEED')
    print(f"Seed is: {SEED}")

    if not SEED:
        sys.exit(1)

    random.seed(int(SEED, 16))

    with open("rockyou.txt", "r") as f:
        lines = f.readlines()
        selection = random.choice(lines).strip()
        print(f"Selection is: {selection}")
        return selection

hash = crypt.crypt(password_from_seed(), crypt.mksalt(crypt.METHOD_SHA512))

base_shadow = f"""root:{hash}:18877::::::
bin:!*:18877::::::
daemon:!*:18877::::::
mail:!*:18877::::::
ftp:!*:18877::::::
http:!*:18877::::::
nobody:!*:18877::::::
dbus:!*:18877::::::
systemd-journal-remote:!*:18877::::::
systemd-network:!*:18877::::::
systemd-oom:!*:18877::::::
systemd-resolve:!*:18877::::::
systemd-timesync:!*:18877::::::
systemd-coredump:!*:18877::::::
uuidd:!*:18877::::::
avahi:!*:18877::::::
brltty:!*:18877::::::
colord:!*:18877::::::
gdm:!*:18877::::::
geoclue:!*:18877::::::
polkitd:!*:18877::::::
rtkit:!*:18877::::::
saned:!*:18877::::::
usbmux:!*:18877::::::
git:!*:18877::::::
nm-openvpn:!*:18883::::::
openvpn:!*:18883::::::
tss:!*:19048::::::
qemu:!*:19128::::::
"""

with open("shadow", "w") as f:
    f.write(base_shadow)