# gpwd

Golang CLI app which generate random password(s) with API :
<https://www.motdepasse.xyz/api>

## Usage

- Generate password :

```
$ gpwd [-l | --length <length>] [-q | --quantity <quantity>] [--no-specials-char] [-e | --export] [-s | statistic]
```
Example : generate 4 password of 15 chars of length with an export
```
$ gpwd -l 15 -q 4 -e
-> 1sp6DbYCk7%#S&4
w37SNc$dh7#I}LN
W2?x8Z7+hfr7Oµz
-2/jJc0.zzPJP[2
[SUCCESS] 4 password(s) export in 'password.txt'
```

- Show version :

```
$ gpwd [-v | --version]
-> gpwd version x.x.x
```

- Helping message

```
$ gpwd [-h | --help]
-> doc ...
```

## TODO

- Add more control to generate password(s)

## License

The Apache License (APACHE) - see [`LICENSE`](./LICENSE) for more details
