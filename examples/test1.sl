init {
  prefix = (pref) {
    return (phrase) { print pref, phrase }
  }
  hi = prefix('Hi')
  hi('Jack')
  goodnight = prefix('Goodnight')
  goodnight('Emily')
  hi('Emily')

  print ''
  arr = ['Jack', 'Emily', 'Buddy']
  arr(hi)

  val = () { 
    val = 34 
  }
  print val()
}
