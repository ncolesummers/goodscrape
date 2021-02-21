# goodscrape
A CLI application to scrape quotes from Goodreads
I would prefer to use their API, but they do not provide new keys 😢

## Usage
`$ goodscrape quotes <query>`

Output is in JSON format.  The filename can be specified with **-f**.

The minimum *amount* of results can be specified with **-a**.

`$ goodscrape quotes -a 100 -f got.json game of thrones`
```json
{
    "Author": "George R.R. Martin,\n  A Game of Thrones",
    "Content": "“... a mind needs books as a sword needs a whetstone, if it is to keep its edge.”"
  },
```
