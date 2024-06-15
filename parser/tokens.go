package parser

type Token int

// syntax tokens
const (
	//
	// meta
	//
	IDENT      = iota
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
	// keywords: constraints
	CONS_PK        // pk (short for primary key)
	CONS_PRIMARY   // primary (followed by key)
	CONS_KEY       // key (with primary prefixed)
	CONS_NULL      // null
	CONS_NOT       // not
	CONS_INCREMENT // increment
	CONS_UNIQUE

	NOTE

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

)
