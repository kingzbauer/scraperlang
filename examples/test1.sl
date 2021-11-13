init {
  print 'Hello, World'
  get 'https://chura.co.ke', {'Content-Type': 'application/json'}
  get 'https://chura.co.ke', {'Content-Type': 'application/json'}
  @help get 'https://google.com', {'Content-Type': 'application/json'}
}

default {
  print 'This is the status', status
}

help {
   print 'in the help closure'
   print status
   get 'https://facebook.com'
}
