package parser

var (
	GrammarProductions = []string{
		"A -> B",
		"B -> C",
		"C -> D",
		"C -> E",
		"D -> {{ F }}",
		"D -> {{ G }}",
		"E -> {{: F }}",
		"H -> ID ( I )",
		"I -> J",
		"I -> ε",
		"J -> J , K",
		"J -> K",
		"K -> VAR_NAME = K",
		"K -> L",
		"L -> UNARY_OP L",
		"L -> M",
		"M -> M LOGIC_OP U",
		"M -> U",
		"U -> U REL_OP N",
		"U -> N",
		"N -> N ADD_OP O",
		"N -> O",
		"O -> O MULT_OP P",
		"O -> P",
		"P -> VAR_NAME",
		"P -> STRING",
		"P -> NUM",
		"P -> ( K )",
		"Q -> IF ( K ) G S T END",
		"S -> ELSE_IF ( K ) G",
		"S -> ε",
		"T -> ELSE G",
		"T -> ε",
		"R -> FOR ( ID IN STRING ) G END",
		"R -> FOR ( ID IN VAR_NAME ) G END",
		"R -> FOR ( ID IN H ) G END",
		"F -> K",
		"F -> H",
		"G -> G F ;",
		"G -> ε",
	}
	LR1ParseTable [][]string = [][]string{
		{"s4", "", "s9", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "130", "1", "2", "3", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r1", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r2", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r3", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r39", "", "s18", "s34", "", "s42", "s50", "", "", "", "s122", "s126", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "5", "7", "12", "", "", "14", "31", "51", "75", "94", "106", "", "", "", "", "70"},
		{"", "s6", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r4", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "s8", "", "s22", "s37", "", "s48", "s56", "", "", "", "s125", "s129", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "16", "", "13", "", "", "15", "41", "57", "90", "97", "109", "", "", "", "", "73"},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r5", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "s18", "s34", "", "s42", "s50", "", "", "", "s122", "s126", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "10", "", "12", "", "", "14", "31", "51", "75", "94", "106", "", "", "", "", "70"},
		{"", "s11", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r6", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r37", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r37", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r36", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r36", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "s17", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r38", "", "r38", "r38", "", "r38", "r38", "", "", "", "r38", "r38", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s19", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "r9", "s44", "s52", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "20", "26", "29", "39", "54", "81", "95", "107", "", "", "", "", "71"},
		{"", "", "", "", "", "s21", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r7", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s23", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "r9", "s44", "s52", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "24", "26", "29", "39", "54", "81", "95", "107", "", "", "", "", "71"},
		{"", "", "", "", "", "s25", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r7", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r8", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "s27", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "", "s44", "s52", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "28", "39", "54", "81", "95", "107", "", "", "", "", "71"},
		{"", "", "", "", "", "r10", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r10", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r11", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r11", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s34", "", "s42", "s50", "", "", "", "s122", "s126", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "43", "31", "51", "75", "94", "106", "", "", "", "", "70"},
		{"", "r13", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "", "s44", "s52", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "47", "39", "54", "81", "95", "107", "", "", "", "", "71"},
		{"", "", "", "", "s36", "", "s45", "s53", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "46", "38", "55", "83", "96", "108", "", "", "", "", "72"},
		{"", "", "", "", "s36", "", "s45", "s53", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "110", "38", "55", "83", "96", "108", "", "", "", "", "72"},
		{"", "", "", "", "s36", "", "s45", "s53", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "112", "38", "55", "83", "96", "108", "", "", "", "", "72"},
		{"", "", "", "", "s36", "", "s45", "s53", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "113", "38", "55", "83", "96", "108", "", "", "", "", "72"},
		{"", "", "", "", "s36", "", "s45", "s53", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "116", "38", "55", "83", "96", "108", "", "", "", "", "72"},
		{"", "", "", "", "", "r13", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r13", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r13", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s37", "", "s48", "s56", "", "", "", "s125", "s129", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "49", "41", "57", "90", "97", "109", "", "", "", "", "73"},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r13", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r24", "", "", "", "", "", "", "r24", "r24", "r24", "", "", "", "", "", "", "", "", "", "", "s30", "r24", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r12", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r24", "", "", "r24", "r24", "r24", "", "", "", "", "", "", "", "", "", "r24", "s32", "r24", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r24", "", "", "r24", "r24", "r24", "", "", "", "", "", "", "", "", "", "", "s33", "r24", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r12", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r12", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r12", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r24", "r24", "r24", "", "", "", "", "", "", "", "", "r24", "", "s40", "r24", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r12", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s34", "", "s118", "s50", "", "", "", "s122", "s126", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "58", "51", "75", "94", "106", "", "", "", "", "70"},
		{"", "r15", "", "", "", "", "", "", "s62", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "", "s119", "s52", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "60", "54", "81", "95", "107", "", "", "", "", "71"},
		{"", "", "", "", "s36", "", "s120", "s53", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "59", "55", "83", "96", "108", "", "", "", "", "72"},
		{"", "", "", "", "", "r15", "", "", "s64", "", "", "", "", "", "", "", "", "", "", "", "r15", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r15", "", "", "s65", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s37", "", "s121", "s56", "", "", "", "s125", "s129", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "61", "57", "90", "97", "109", "", "", "", "", "73"},
		{"", "", "", "", "", "", "", "", "s68", "", "", "", "", "", "", "", "", "", "", "r15", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r14", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r14", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r14", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r14", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "r14", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s34", "", "s118", "", "", "", "", "s122", "s126", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "75", "94", "106", "", "", "", "", "63"},
		{"", "r16", "", "", "", "", "", "", "r16", "", "", "", "", "", "", "", "", "", "", "", "", "", "s74", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "", "s119", "", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "81", "95", "107", "", "", "", "", "66"},
		{"", "", "", "", "s36", "", "s120", "", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "83", "96", "108", "", "", "", "", "67"},
		{"", "", "", "", "", "r16", "", "", "r16", "", "", "", "", "", "", "", "", "", "", "", "r16", "", "s79", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r16", "", "", "r16", "", "", "", "", "", "", "", "", "", "", "", "", "", "s80", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s37", "", "s121", "", "", "", "", "s125", "s129", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "90", "97", "109", "", "", "", "", "69"},
		{"", "", "", "", "", "", "", "", "r16", "", "", "", "", "", "", "", "", "", "", "r16", "", "", "s89", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r17", "", "", "", "", "", "", "r17", "", "", "", "", "", "", "", "", "", "", "", "", "", "s74", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r17", "", "", "r17", "", "", "", "", "", "", "", "", "", "", "", "r17", "", "s79", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r17", "", "", "r17", "", "", "", "", "", "", "", "", "", "", "", "", "", "s80", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r17", "", "", "", "", "", "", "", "", "", "", "r17", "", "", "s89", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s34", "", "s118", "", "", "", "", "s122", "s126", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "76", "94", "106", "", "", "", "", ""},
		{"", "r19", "", "", "", "", "", "", "r19", "s77", "", "", "", "", "", "", "", "", "", "", "", "", "r19", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r18", "", "", "", "", "", "", "r18", "s77", "", "", "", "", "", "", "", "", "", "", "", "", "r18", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s34", "", "s118", "", "", "", "", "s122", "s126", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "78", "106", "", "", "", "", ""},
		{"", "r20", "", "", "", "", "", "", "r20", "r20", "s98", "", "", "", "", "", "", "", "", "", "", "", "r20", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "", "s119", "", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "82", "95", "107", "", "", "", "", ""},
		{"", "", "", "", "s36", "", "s120", "", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "84", "96", "108", "", "", "", "", ""},
		{"", "", "", "", "", "r19", "", "", "r19", "s85", "", "", "", "", "", "", "", "", "", "", "r19", "", "r19", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r18", "", "", "r18", "s85", "", "", "", "", "", "", "", "", "", "", "r18", "", "r18", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r19", "", "", "r19", "s86", "", "", "", "", "", "", "", "", "", "", "", "", "r19", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r18", "", "", "r18", "s86", "", "", "", "", "", "", "", "", "", "", "", "", "r18", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "", "s119", "", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "87", "107", "", "", "", "", ""},
		{"", "", "", "", "s36", "", "s120", "", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "88", "108", "", "", "", "", ""},
		{"", "", "", "", "", "r20", "", "", "r20", "r20", "s100", "", "", "", "", "", "", "", "", "", "r20", "", "r20", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r20", "", "", "r20", "r20", "s101", "", "", "", "", "", "", "", "", "", "", "", "r20", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s37", "", "s121", "", "", "", "", "s125", "s129", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "91", "97", "109", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r19", "s92", "", "", "", "", "", "", "", "", "", "r19", "", "", "r19", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r18", "s92", "", "", "", "", "", "", "", "", "", "r18", "", "", "r18", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s37", "", "s121", "", "", "", "", "s125", "s129", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "93", "109", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r20", "r20", "s104", "", "", "", "", "", "", "", "", "r20", "", "", "r20", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r21", "", "", "", "", "", "", "r21", "r21", "s98", "", "", "", "", "", "", "", "", "", "", "", "r21", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r21", "", "", "r21", "r21", "s100", "", "", "", "", "", "", "", "", "", "r21", "", "r21", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r21", "", "", "r21", "r21", "s101", "", "", "", "", "", "", "", "", "", "", "", "r21", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r21", "r21", "s104", "", "", "", "", "", "", "", "", "r21", "", "", "r21", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s34", "", "s118", "", "", "", "", "s122", "s126", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "99", "", "", "", "", ""},
		{"", "r22", "", "", "", "", "", "", "r22", "r22", "r22", "", "", "", "", "", "", "", "", "", "", "", "r22", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s35", "", "s119", "", "", "", "", "s123", "s127", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "102", "", "", "", "", ""},
		{"", "", "", "", "s36", "", "s120", "", "", "", "", "s124", "s128", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "103", "", "", "", "", ""},
		{"", "", "", "", "", "r22", "", "", "r22", "r22", "r22", "", "", "", "", "", "", "", "", "", "r22", "", "r22", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r22", "", "", "r22", "r22", "r22", "", "", "", "", "", "", "", "", "", "", "", "r22", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "s37", "", "s121", "", "", "", "", "s125", "s129", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "105", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r22", "r22", "r22", "", "", "", "", "", "", "", "", "r22", "", "", "r22", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r23", "", "", "", "", "", "", "r23", "r23", "r23", "", "", "", "", "", "", "", "", "", "", "", "r23", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r23", "", "", "r23", "r23", "r23", "", "", "", "", "", "", "", "", "", "r23", "", "r23", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r23", "", "", "r23", "r23", "r23", "", "", "", "", "", "", "", "", "", "", "", "r23", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r23", "r23", "r23", "", "", "", "", "", "", "", "", "r23", "", "", "r23", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "s111", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r27", "", "", "", "", "", "", "r27", "r27", "r27", "", "", "", "", "", "", "", "", "", "", "", "r27", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "s114", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "s115", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r27", "", "", "r27", "r27", "r27", "", "", "", "", "", "", "", "", "", "r27", "", "r27", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r27", "", "", "r27", "r27", "r27", "", "", "", "", "", "", "", "", "", "", "", "r27", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "s117", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r27", "r27", "r27", "", "", "", "", "", "", "", "", "r27", "", "", "r27", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r24", "", "", "", "", "", "", "r24", "r24", "r24", "", "", "", "", "", "", "", "", "", "", "", "r24", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r24", "", "", "r24", "r24", "r24", "", "", "", "", "", "", "", "", "", "r24", "", "r24", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r24", "", "", "r24", "r24", "r24", "", "", "", "", "", "", "", "", "", "", "", "r24", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r24", "r24", "r24", "", "", "", "", "", "", "", "", "r24", "", "", "r24", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r25", "", "", "", "", "", "", "r25", "r25", "r25", "", "", "", "", "", "", "", "", "", "", "", "r25", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r25", "", "", "r25", "r25", "r25", "", "", "", "", "", "", "", "", "", "r25", "", "r25", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r25", "", "", "r25", "r25", "r25", "", "", "", "", "", "", "", "", "", "", "", "r25", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r25", "r25", "r25", "", "", "", "", "", "", "", "", "r25", "", "", "r25", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "r26", "", "", "", "", "", "", "r26", "r26", "r26", "", "", "", "", "", "", "", "", "", "", "", "r26", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r26", "", "", "r26", "r26", "r26", "", "", "", "", "", "", "", "", "", "r26", "", "r26", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "r26", "", "", "r26", "r26", "r26", "", "", "", "", "", "", "", "", "", "", "", "r26", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "r26", "r26", "r26", "", "", "", "", "", "", "", "", "r26", "", "", "r26", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "acct", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	}

	SymbolColMapping = map[string]int{
		"{{":       0,
		"}}":       1,
		"{{:":      2,
		"ID":       3,
		"(":        4,
		")":        5,
		"VAR_NAME": 6,
		"UNARY_OP": 7,
		"LOGIC_OP": 8,
		"ADD_OP":   9,
		"MULT_OP":  10,
		"STRING":   11,
		"NUM":      12,
		"IF":       13,
		"ELSE_IF":  14,
		"ELSE":     15,
		"END":      16,
		"FOR":      17,
		"IN":       18,
		";":        19,
		",":        20,
		"=":        21,
		"REL_OP":   22,
		"$":        23,
		"B":        24,
		"C":        25,
		"D":        26,
		"E":        27,
		"F":        28,
		"G":        29,
		"H":        30,
		"I":        31,
		"J":        32,
		"K":        33,
		"L":        34,
		"M":        35,
		"N":        36,
		"O":        37,
		"P":        38,
		"Q":        39,
		"R":        40,
		"S":        41,
		"T":        42,
	}
)