# Install:
1.
```bash
npm install -g retire
```
2.
```bash
go install -v github.com/whalebone7/gretire@latest
```


# Usage
This tool allows you to scan multiple javascript endpoints by embedding retire.js in go lang:

```bash
cat javascriptUrls.txt | gretire
```
Or 

```bash
echo "https://example.com/assets/tracking.js" | gretire
```
