package tokens

type Token uint64

var EOFChar = rune(0)

// syntax tokens
const (
	//
	// meta
	//
	IDENT      = 1 << iota
	WHITESPACE // whitespace
	LINEBR     // line break \n
	EOF        // end-of-file
	ILLEGAL    // illegal token
	UNKOWN     // unkown token

	//
	// keywords
	//
	TABLE              // Table
	ENUM               // enum
	COL_SETTING_CUSTOM // something like id "bigint unsigned" [pk]
	REF_CAP            // Ref
	REF_LOW            // ref (inline)
	// keywords: constraints
	CONS_PK        // pk (short for primary key)
	CONS_PRIMARY   // primary (followed by key)
	CONS_KEY       // key (with primary prefixed)
	CONS_NULL      // null
	CONS_NOT       // not
	CONS_INCREMENT // increment
	CONS_UNIQUE

	NOTE

	REL_1T1 // -
	REL_1TM // <
	REL_MT1 // >
	REL_MTN // <>

	//
	// delimiters
	//
	SLASH        // /
	BRACE_OPEN   // {
	BRACE_CLOSE  // }
	SQUARE_OPEN  // [
	SQUARE_CLOSE // ]
	ROUND_OPEN   // (
	ROUND_CLOSE  // )
	COLON        // :
	COMMA        // ,
	APOSTROPHE   // '
	// BACKTICK     // `
	QUOTATION // \"
	DOT       // .

	//
	// GROUPS
	//
	G_RELATION_TYPE = REL_1T1 | REL_MT1 | REL_1TM | REL_MTN
)

func MapLiteral(literal string) Token {
	switch literal {
	case "Table":
		return TABLE
	case "enum":
		return ENUM
	case "pk":
		return CONS_PK
	case "primary":
		return CONS_PRIMARY
	case "key":
		return CONS_KEY
	case "null":
		return CONS_NULL
	case "not":
		return CONS_NOT
	case "increment":
		return CONS_INCREMENT
	case "unique":
		return CONS_UNIQUE
	case "note":
		return NOTE
	case "Ref":
		return REF_CAP
	case "ref":
		return REF_LOW
	case "<>":
		return REL_MTN
	default:

	}
	return IDENT
}

func MapChar(char rune) Token {
	switch char {
	case '\n':
		return LINEBR
	case EOFChar:
		return EOF
	case '/':
		return SLASH
	case '{':
		return BRACE_OPEN
	case '}':
		return BRACE_CLOSE
	case '[':
		return SQUARE_OPEN
	case ']':
		return SQUARE_CLOSE
	case '"':
		return QUOTATION
	case ',':
		return COMMA
	case ':':
		return COLON
	case '.':
		return DOT
	case '-':
		return REL_1T1
	case '>':
		return REL_MT1
	case '<':
		return REL_1TM
	default:
	}
	return UNKOWN
}
