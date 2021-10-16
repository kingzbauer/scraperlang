package scraperlang

/*
	document				-> global_defs* ;
	global_defs 		-> NEWLINE* tagged_closure ( NEWLINE+ | EOF ) ;
	tagged_closure	-> IDENT "{" ( NEWLINE+ expr_statements* )? "}";
	builtin_funcs		-> getExpr | printExpr ;
	expr_statements	-> ( expressions | builtin_funcs ) terminator ;
	expression 			-> assignExpr |
										 callExpr |
										 closureExpr |
										 accessExpr |
										 htmlAttrAccessor |
										 arrayExpr |
										 mapExpr |
										 mapAccess |
										 arrayAccess |
										 primary ;
	getExpr					-> tag? "get" expression ("," expression) ;
	tag							-> "@"IDENT ;
	assignExpr 			-> IDENT "=" ( primary | expression ) ;
	printExpr				-> expression ( "," expression )* ;
	callExpr				-> IDENT expression ( "," expression ) ;
	closureExpr			-> "(" params? ")" "{" expr_statements* "}" ;
	accessExpr 			-> IDENT "." IDENT ;
	htmlAttrAccessor-> IDENT "~" IDENT ;
	arrayExpr				-> "[" NEWLINE* arrayEntry NEWLINE* ( "," NEWLINE* valueExpr NEWLINE )* "]" ;
	valueExpr				-> ( primary
										 | callExpr
										 | arrayExpr
										 | mapExpr
										 | htmlAttrAccessor
										 | accessExpr ) ;
	mapExpr					-> "{" NEWLINE* mapEntry NEWLINE* ( "," NEWLINE* mapEntry NEWLINE* )* "}" ;
	mapEntry				-> STRING ":" valueExpr ;
	mapAccess				-> IDENT "[" STRING "]" ;
	primary					-> STRING | NUMBER | TRUE | FALSE | NIL | IDENT ;
	terminator			-> NEWLINE+ ;
*/
