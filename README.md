# gpwd

Golang CLI app which generate random password(s) with API :
<https://www.motdepasse.xyz/api>

## Description

This CLI app allows to generate passwords or keys with a limit of 512 characters in length and a maximum of 30 passwords at the same time. With the possibility of making exports and having an eye on statistics

## Usage

- Generate password :

```
$ gpwd [-l | --length <length>] [-q | --quantity <quantity>] [--no-specials-char] [-e | --export] [-s | --statistic]
```
Example : generate 4 password of 15 chars of length with an export and statistic speed log
```
$ gpwd -l 15 -q 4 -es
-> 1sp6DbYCk7%#S&4
w37SNc$dh7#I}LN
W2?x8Z7+hfr7Oµz
-2/jJc0.zzPJP[2
[SUCCESS] 4 password(s) export in 'password.txt'
Finished in 100ms
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

## License

The Apache License (APACHE) - see [`LICENSE`](./LICENSE) for more details
