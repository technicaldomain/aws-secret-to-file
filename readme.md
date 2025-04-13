# AWS Secret to file command line utility

Simple command line utility to retrieve secret from AWS Secret Manager and store it as file


### Limitations

Only plain text and binary secrets are supported now


### How to use

Retrieve single secret

```/bin/aws-secret-to-file --secret=/secret/name/here --output=./location/for/the/file```

Retrieve binary secret

```/bin/aws-secret-to-file --secret=/secret/name/here --output=./location/for/the/file --binary```

Retrieve multiple secrets

```sh
/bin/aws-secret-to-file \
    --secret=/secret/name/here1 --output=./location/for/the/file1 \
    --secret=/secret/name/here2 --output=./location/for/the/file2
```

Create binary secret from command line

```sh
aws secretsmanager create-secret \
  --name /binary/secret \
  --secret-binary fileb://foo.bar
```


## License

Apache 2 Licensed. For more information please see [LICENSE](LICENSE)
