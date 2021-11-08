init {
  prefix = (pref) {
    fn = (phrase) { print pref, phrase }
  }
  hi = prefix('Hi')
  hi('Jack')
  goodnight = prefix('Goodnight')
  goodnight('Emily')
  hi('Emily')

  print ''
  arr = ['Jack', 'Emily', 'Buddy']
  arr(hi)

  val = () { print "jack" }
  print val
  val()
}
