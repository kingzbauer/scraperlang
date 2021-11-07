init {
  a = [2, 3, 4, 5, (v) {
    name = {"location": v}
  }]
  print a[4]("Ruiru")['location']
}
