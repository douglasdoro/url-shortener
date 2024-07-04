# url-shortener
URL shortening service developed in the book <a href="http://www.casadocodigo.com.br/products/livro-google-go" rel="nofollow">"Programando em Go: crie aplicações com a linguagem do Google"</a>, <a href="http://www.casadocodigo.com.br" rel="nofollow">Casa do Código</a>.

## How to

```bash
# Run
go run .

# build
go build -o urlshortener

# to see the options
./urlshortener -h

# Example
./urlshortener -l=false -p=4000

```

### Request samples
```bash
curl --location 'http://localhost:3000/api/tiny' \
--header 'Content-Type: text/plain' \
--data 'https://www.terra.com.br'

curl --location 'http://localhost:3000/api/stats/1US3T'

curl --location 'http://localhost:3000/r/1US3T'
```

