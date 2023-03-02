package inputinfo

// func parseVideoLine(line string) *videostream {

// 	reader := parser.NewReader(line)
// 	fn := vidStatements()
// 	res, err := parser.Run(line, fn)
// 	fmt.Println("+++++")
// 	fmt.Println(res.ToString())
// 	fmt.Println(err)
// 	fmt.Println("+++++")
// 	// fn := parser.JOIN(

// 	// 	"ab", parser.Keep("middle", "cd"), "ef",
// 	// )

// 	reader.Read()
// 	return nil
// }

// func vidStatements() parser.ParserFunc {
// 	// fn := parser.JOIN(
// 	// 	parser.OPT(" "),
// 	// )
// 	fn1 := parser.Keep("stream", parser.JOIN(
// 		parser.OPT(
// 			parser.WHILE(" "),
// 		),
// 		parser.JOIN("Stream"),
// 		parser.OPT(
// 			parser.WHILE(" "),
// 		),
// 		parser.JOIN(
// 			"#",
// 			parser.LIT(parser.Digit),
// 			":",
// 			parser.LIT(parser.Digit),
// 		),
// 		parser.OPT(
// 			parser.JOIN("(", parser.WHILE(parser.Ident()))),
// 	),
// 	)

// 	return fn1
// }
