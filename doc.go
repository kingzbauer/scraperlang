package scraperlang

/*
	document				-> global_defs* EOF;
	global_defs 		-> tagged_closure ;
	tagged_closure	-> IDENT "{" expr_statements* "}" ;
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
	closureExpr			-> "{" params? expr_statements* "}" ;
	accessExpr 			-> IDENT "." IDENT ;
	htmlAttrAccessor-> IDENT "~" IDENT ;
	arrayExpr				-> "[" NEWLINE* arrayEntry NEWLINE* ( "," NEWLINE* arrayEntry NEWLINE* )* "]" ;
	arrayEntry			-> ( primary
										 | callExpr
										 | arrayExpr
										 | mapExpr
										 | htmlAttrAccessor
										 | accessExpr ) ;
	mapExpr					-> "{" NEWLINE* mapEntry NEWLINE* ( "," NEWLINE* mapEntry NEWLINE* )* "}" ;
	mapEntry				-> STRING ":" arrayEntry ;
	mapAccess				-> IDENT "[" STRING "]" ;
	primary					-> STRING | NUMBER | TRUE | FALSE | NIL | IDENT ;
	terminator			-> NEWLINE+ ;
*/
