# Bitcoin HD wallet
- [Objective](#objective)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Demo](#demo)
- [Security](#security)

## Objective
To create an HTTP API server that supports the following operations: 
1. Generate a random mnemonic words following BIP39 standard
2. Generate a Hierarchical Deterministic (HD) Segregated Witness (SegWit) bitcoin address from a given seed and path
3. Generate an n-out-of-m Multisignature (multi-sig) Pay-To-Script-Hash (P2SH) bitcoin address, where n, m, and addresses can be specified 

## Prerequisites
[go](https://golang.org/) 1.17 or newer.

## Installation
Download dependencies
```
make dep
```    
To run server
```
make run
```
To run test cases
```
make test
```
Show code coverage
```
make coverage
```

## API documentation
To view detail documentation on Api endpoint. [Click here to view swagger](https://leuhk.github.io/bitcoin-hdwallet/)

## Demo
### 1. Generate a random mnemonic words following BIP39 standard.

```
curl --location --request GET 'http://localhost:8080/util/new-mnemonic' \
--data-raw ''
```
  Exmaple response
```
{
"BIP39Mnemonic":"run tiny elevator venture nice hazard about price begin steel long penalty check person scissors sister odor nest ribbon issue second quarter grow vendor",
"BIP39Seed":"76a26fce9718552dcdadabb87cb93c7aa6665baf70c95c44a422614e1fc3cdd875efa211f1f051a753adb32fd4e8e081157a3549e9b89f7095a2e59802109f61"
}
```
**please note:**
 - Program is set to only generated mnemonic with 24 words.
 - The generated BIP39 seed is not protected with a passpharse 
### 2. Generate a Hierarchical Deterministic (HD) Segregated Witness (SegWit) bitcoin address from a given seed and path
```
curl --location --request POST 'http://localhost:8080/util/hd-wallet' \
--header 'Content-Type: application/json' \
--data-raw '{
    "seed":"3fdf3c7c40ef678dd8950caac27f8006e27fdfab5e379ff7e3cef34a0226df830a49ca85476e7873e096ca6127d365995f6f135c71c27e8efe6cd1c497f6003f",
    "path":"m / 44'\'' / 0'\'' / 0'\'' / 0 / 0"
}'
```
Exmaple response
```
{
    "WIF": "L1dbCB2GPDDzxNJJ7s7h7a84NSXViBNwJeLWTS2gJMtdapfkhmg8",
    "bip32ExtendedPrivateKey": "xprvA319vcXCEmKZe8ens22j5pGfY6WR2yEeBnPPMPE2CDpn4JaoXyYjsHWyDeDbXFXDWwuJAgbJve2772PRfVrY6jFUBj43JDbXMJ5EZQYKDhM",
    "bip32ExtendedPublicKey": "xpub6FzWL84658srrcjFy3ZjSxDQ68LuSRxVZ1Jz9mddkZMkw6ux5WrzR5qT4wSsnG7zpfQFrAeQDeoRzec8xXy5FRz8ZDewDG3NV8nDFNjYrjZ",
    "bip32RootKey": "xprv9s21ZrQH143K2pnPh3AEko6pZTqmoyFW3Kt8heSpwhSzSfSJP3T9rFome7xNkvk9GaW7M91QEvkbP22z6HwhvFqTtuisH5hHPTu5xDBQRkG",
    "p2pkhAddress": "1AKdhnB63swG2XpSuuXWP8MqP596zRJJCz",
    "segwitBech32": "bc1qvcl5rm6vz7eaxed9zrcftax0ywsc29y8zgj876",
    "segwitNested": "3MTm6vsDYfyQCSTeex9ZHWYd9zkUetUn7c"
}
```
**please note:**
  - the path format must follow the [BIP44 standard](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#examples) where a harden value must have a **'** suffix
  ```
  m / purpose' / coin_type' / account' / change / address_index
  ```
  - purpose, coin_type and account must be a harden value, whereas change and address_index must not be harden
  - Keys are generated for mainnet enviornemnt 
  - Segwit address consits of two types:
    - 'bc1' prefixed Native Segwit (bech-32)
    - '3' prefixed nested segwit (p2wpkh-p2sh)
### 3.Generate an n-out-of-m Multisignature (multi-sig) Pay-To-Script-Hash (P2SH) bitcoin address, where n, m, and addresses can be specified

```
curl --location --request POST 'http://localhost:8080/util/multi-sig-p2sh' \
--header 'Content-Type: application/json' \
--data-raw '{
    "n":1,
    "m":2,
    "wif":["cQpHXfs91s5eR9PWXui6qo2xjoJb2X3VdUspwKXe4A8Dybvut2rL","cVgxEkRBtnfvd41ssd4PCsiemahAHidFrLWYoDBMNojUeME8dojZ"]
}'

```
Exmaple response
```
{
    "address": "3MqSiHLbK6M8YUL8sXKiULeiRSvckJV74h"
}
```
**please note:**

m = the minimum number of signatures that are required.

n = the total number of public keys, used in multi-sig script.

m and n can be up to 16

wif = private keys in WIF format

## Security 
for a good security practice, HD wallet should generate a seed phrase with length of words greater than or equal to 12 to improve security, as it will be more complex to derive the generated address

In an ideal production grade application, the seed should use some sort of encryption standard before sending to both client side and server side. This will avoid man-in-the-middle attacks sniffing Https Request.
The database should also not store seed in a raw format.
    