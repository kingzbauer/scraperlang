package scraperlang

/*
	document				-> global_defs* ;
	global_defs 		-> NEWLINE* tagged_closure* ( NEWLINE+ | EOF ) ;
	tagged_closure	-> IDENT body ;
	body 						-> "{" ( NEWLINE+ expr_statements* )? "}" ;
	builtin_funcs		-> ( getExpr | printExpr ) NEWLINE ;
	expr_statements	-> ( assign | builtin_funcs | callExpr ) NEWLINE ;
	getExpr					-> tag? "get" expression ("," expression) ;
	tag							-> "@"IDENT ;
	printExpr				-> "print" expression ( "," expression )* ;
	closure					-> "(" params? ")" body ;
	arrayExpr				-> "[" NEWLINE* expression NEWLINE* ( "," NEWLINE* expression NEWLINE )* "]" ;
	mapExpr					-> "{" NEWLINE* mapEntry NEWLINE* ( "," NEWLINE* mapEntry NEWLINE* )* "}" ;
	mapEntry				-> STRING ":" expression ;

	callExpr				-> IDENT expression ( "," expression ) ;
	assign					-> IDENT "=" ( expression ) ;
	expression 			-> htmlAttrAccessor ;
	htmlAttrAccessor-> accessor ( "~" IDENT )? ;
	accessor 				-> ( ( primary ( ( "(" arguments? ")" ) |
										 ( "[" expression "]" ) |
									 		"." IDENT )* ) | mapExpr | arrayExpr | closure ) ;
	primary					-> STRING | NUMBER | TRUE | FALSE | NIL | IDENT ;
*/
