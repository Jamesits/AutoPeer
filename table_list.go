package main

func addTable(conf *config, table string) {
	if len(table) == 0 {
		return
	}
	if ok, _ := inArray(table, conf.tables); !ok {
		conf.tables = append(conf.tables, table)
	}
}
