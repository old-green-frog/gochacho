package server

type TemplateOption struct {
	Name string
	Key  string
	Date bool
}

type TemplateData struct {
	Options []TemplateOption
}

var ROLES_OPTIONS = map[string][]TemplateOption{
	"mech": []TemplateOption{
		TemplateOption{"Список выполненных ремонтов", "work", true},
		TemplateOption{"Список не выполненных ремонтов", "not_work", false},
		TemplateOption{"Дефектная ведомость", "defectlist", false},
	},
	"main_ch": []TemplateOption{
		TemplateOption{"Сводный отчет о выполненных ремонтах", "report", true},
		TemplateOption{"Сведения о поверках счётчиков", "check", true},
	},
	"caseer": []TemplateOption{
		TemplateOption{"Список зарегестрированных счетчиков", "meters", false},
		TemplateOption{"Список выполненных поверок счётчиков", "checks", true},
		TemplateOption{"Список установленных счётчиков", "installed", true},
		TemplateOption{"Список необходимых поверок счётчиков", "will_checks", false},
	},
}
