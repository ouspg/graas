# Generate a shadow file for Linux system for password cracking purposes

https://man7.org/linux/man-pages/man5/shadow.5.html

Shadow compatible hash with included salt and option `$6` (SHA-512 hashing):

```console
python3 -c 'import crypt; print(crypt.crypt("password", crypt.mksalt(crypt.METHOD_SHA512)))'
```

Produces for example:

> `$6$WSE04aebjazYHra.$F1Ei6ujDlBUCRvBGFGeUn3iuzYihzr6J57i5gvn5ZPJeJu91vKOlhcLGEOPhuo2JtMaAfmYqwuhJb8sKRKiz7.`

## Useful

https://unix.stackexchange.com/questions/430141/how-to-find-the-hashing-algorithm-used-to-hash-passwords

https://serverfault.com/questions/330069/how-to-create-an-sha-512-hashed-password-for-shadow
