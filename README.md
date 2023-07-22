# Encrypted Go Backup

This project uses a user provided key to encrypt all files in a given file path and pushes them to a user provided AWS S3 bucket.

## Generating the Key

This project expects an RSA key.

`ssh-keygen -t rsa -b 4096 -f <key file name>`

Do NOT use a password on the key

Make sure the public key and private key share the same name, with the private key ending in `.pub` and no extension on the private key

## Generating the config.json

This project runs off of a user provided Config file named `config.json`

| Config item key | Description | required | default |
|-|-|-|-|
|`s3.bucket`| S3 bucket in the AWS Account to push to| X |None|
|`s3.prefix`| Prefix to append to the file names when pushing to S3 | |None|
|`key.fileName`|the filename of the public key/private key combo. Only provide the base name. the `.pub` will be appended|X|None|
|`key.path`|The full path to the location of the key files||`~/.ssh/`*|
|`backupPath`|The path to find, encrypt, and upload files from. (must be directory)|X|None|
|`decryptPath`|The path to store decrypted files on a decrypt run||None|

### Example config.json

```json
{
    "s3": {
        "bucket": "bucket-name-to-use",
        "prefix": "prefix-to-append"
    },
    "key": {
        "fileName": "key-name-from-above",
        "path": "key-location-on-computer"
    },
    "backupPath": "folder-location-to-start-backup",
    "decryptPath": "folder-to-store-decrypted-files"
}
```

## Building and Running

Run `go build` and a binary executable `backup` will be generated. Run this manually or on a schedule with arguments provided in the arguments section

\* - **Do NOT run as root. Aside from the glaring security concerns, file paths for SSH keys will not be the same unless manually specified**

## Arguements

By default, the script will ENCRYPT files without any argument present. The only supported arguements are `encrypt` and `decrypt`.

### Example Usage

`./backup`\
`./backup encrypt`\
`./backup decrypt`