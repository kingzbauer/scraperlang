init {
  get "https://leannx.com", { "Authorization": "Token asdflkjasdflj" }
  headers = {
    "Accept": "application/json"
  }
  @anchors get "https://leannx.com", headers

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
  anchors.loop {
    @anchor get it~href
    print it~id
  }
  closure = { item, index ->
    @anchor get item~href
  }
  anchors.loop closure
}
