package core

// EnToUa maps English keyboard layout runes to Ukrainian layout runes
var EnToUa = map[rune]rune{
	'q':'й','w':'ц','e':'у','r':'к','t':'е','y':'н','u':'г','i':'ш','o':'щ','p':'з','[':'х',']':'ї',
	'a':'ф','s':'і','d':'в','f':'а','g':'п','h':'р','j':'о','k':'л','l':'д',';':'ж','\'':'є',
	'z':'я','x':'ч','c':'с','v':'м','b':'и','n':'т','m':'ь',',':'б','.':'ю','/':'.',

	'Q':'Й','W':'Ц','E':'У','R':'К','T':'Е','Y':'Н','U':'Г','I':'Ш','O':'Щ','P':'З',
	'A':'Ф','S':'І','D':'В','F':'А','G':'П','H':'Р','J':'О','K':'Л','L':'Д',
	'Z':'Я','X':'Ч','C':'С','V':'М','B':'И','N':'Т','M':'Ь',
}
// Reverse inverts a map[rune]rune
func Reverse(src map[rune]rune) map[rune]rune {
	// Create a new map with inverted key-value pairs
	dst := make(map[rune]rune, len(src))
	// Iterate over the source map and swap keys and values
	for k, v := range src {
		dst[v] = k
	}
	return dst
}

var UaToEn = Reverse(EnToUa)
