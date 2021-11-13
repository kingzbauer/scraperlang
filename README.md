# Scraper-lang (WIP)
A DSL aimed at making writing web scrapers/crawlers a breeze

## Example

```
init {
  get "https://google.com", { "Authorization": "Token asdflkjasdflj" }
  headers = {
    "Accept": "application/json"
  }
  @anchors get "https://google.com", headers

  urls = ['https://google.com', 'https://chura.co.ke']
  @anchors get urls, headers
}

default {
  print content
}

anchors {
  print headers, headers['Content-Length']
  print headers['Content-Type']
  anchors = jq 'a#base'
  anchors.loop(() {
    @anchor get it~href
    print it~id
  })
  closure = (item, index ) {
    @anchor get item~href
  }
  anchors.loop  closure
}
```

## TODO:

- Make "post" expressions keyword expressions
- Handle quering JSON responses
- Handle GraphQL queries
- Implement Resolver
  - Will check that there are no cycles
  - Check return statement is used properly
  - Check for missing but referenced tagged closures
- Design how to handle terminating conditions for tagged closures that might lead into an infinite loop
