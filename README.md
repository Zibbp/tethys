# Tethys

Simple file upload and download using curl. Files are deleted after one hour. Designed to be deployed internally, not exposed publicly. Tethys does not offer any security features.

Upload a file with:

```sh
curl --upload-file demo.json https://example.com
```

Download a file with:

```sh
// Upload will return the file URL
curl https://example.com/4jfh2/demo.json -o demo.json
```
