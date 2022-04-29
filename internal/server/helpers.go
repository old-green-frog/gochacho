package server

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type TemplateOption struct {
	Name   string
	Key    string
	Date   bool
	Create bool
}

type TemplateData struct {
	Options []TemplateOption
}

var ROLES_OPTIONS = map[string][]TemplateOption{
	"mech": []TemplateOption{
		TemplateOption{"Список выполненных ремонтов", "work", true, true},
		TemplateOption{"Список не выполненных ремонтов", "not_work", false, true},
		TemplateOption{"Дефектная ведомость", "defectlist", false, true},
	},
	"main_ch": []TemplateOption{
		TemplateOption{"Сводный отчет о выполненных ремонтах", "report", true, false},
		TemplateOption{"Сведения о поверках счётчиков", "check", true, false},
	},
	"caseer": []TemplateOption{
		TemplateOption{"Список зарегестрированных счетчиков", "meters", false, false},
		TemplateOption{"Список выполненных поверок счётчиков", "checks", true, true},
		TemplateOption{"Список установленных счётчиков", "installed", true, true},
		TemplateOption{"Список необходимых поверок счётчиков", "will_check", false, false},
	},
}

var QUERIES = map[string]string{
	"work": `SELECT w.typename, e.id as serialnumber, e.brand, br.brigadier_name, d.id, d.date_of_start::varchar, d.desired_time, d.date_of_accept::varchar, EXTRACT( DAY FROM d.date_of_end::timestamp - d.date_of_start::timestamp) as actualtime
			 	FROM defectivelist as d, equipment as e, worktype as w, brigade as br
			 	WHERE d.equipment_id = e.id AND d.worktype_id = w.id AND d.brigade_id = br.id AND d.date_of_end IS NOT NULL
			 	AND d.date_of_start BETWEEN '%s' AND '%s';`,
	"not_work": `SELECT w.typename, e.id, e.brand, d.date_of_start::varchar, ((d.date_of_start + interval '1' day * d.desired_time)::date)::varchar, br.brigadier_name, d.id
					FROM defectivelist as d, equipment as e, worktype as w, brigade as br
					WHERE d.equipment_id = e.id AND d.worktype_id = w.id AND d.brigade_id = br.id AND d.date_of_end IS NULL
					AND d.date_of_start = '%s';`,
	"defectlist": `SELECT e.id, e.brand, '-', d.date_of_start::varchar, d.date_of_end::varchar, EXTRACT( DAY FROM d.date_of_end::timestamp - d.date_of_start::timestamp)
					FROM equipment as e, defectivelist as d
					WHERE e.id = d.equipment_id AND d.id = '%s';`,
	"meters": `SELECT a.id, m.id, a.customer_name, a.customer_address, a.customer_number
					FROM account as a, meter as m
					WHERE a.meter_id = m.id AND m.date_of_plug = '%s';`,
	"check": `SELECT m.date_of_plug::varchar, COUNT(m.id), m.price_for_plug::varchar, m.date_of_check::varchar, m.price_for_check::varchar
				FROM meter as m
				WHERE m.price_for_check IS NOT NULL AND m.date_of_plug BETWEEN '%s' AND '%s'
				GROUP BY m.date_of_plug, m.price_for_plug, m.date_of_check, m.price_for_check;`,
	"report": `SELECT w.typename, COUNT(d.id), d.desired_time, EXTRACT( DAY FROM d.date_of_end::timestamp - d.date_of_start::timestamp) as fact_time, 
				CASE
					WHEN d.desired_time - EXTRACT( DAY FROM d.date_of_end::timestamp - d.date_of_start::timestamp) > 0 THEN d.desired_time - EXTRACT( DAY FROM d.date_of_end::timestamp - d.date_of_start::timestamp)
					ELSE 0
				END
				FROM worktype as w, defectivelist as d
				WHERE w.id = d.worktype_id AND d.date_of_end BETWEEN '%s' AND '%s'
				GROUP BY w.id, d.desired_time, d.date_of_end, d.date_of_start;`,
	"checks": `SELECT a.id, a.customer_name, a.customer_address, m.id, m.date_of_check::varchar, m.price_for_check::varchar
					FROM account as a, meter as m
					WHERE a.meter_id = m.id AND m.price_for_check IS NOT NULL AND m.date_of_check BETWEEN '%s' AND '%s';`,
	"installed": `SELECT a.id, a.customer_name, a.customer_address, m.id, a.customer_number, m.date_of_plug::varchar, m.price_for_plug::varchar
	FROM account as a, meter as m
	WHERE a.meter_id = m.id AND m.date_of_plug BETWEEN '%s' AND '%s';`,
	"will_check": `SELECT a.id, a.customer_name, a.customer_address, m.id, a.customer_number, m.date_of_check::varchar
	FROM account as a, meter as m
	WHERE a.meter_id = m.id AND m.price_for_check IS NULL AND m.date_of_check = '%s';`,
}

var TITLES = map[string][]string{
	"work": {"Вид ремонта",
		"Региональный №",
		"Наименование оборудования, сооружения или строения",
		"Ответственный исполнитель",
		"№ деф. ведом",
		"Дата начала ремонта",
		"Кол-во дней по плану",
		"Дата приемки объектов после ремонта",
		"Кол-во дней факт"},
	"not_work": {"Вид ремонта",
		"Регистрационный номер",
		"Наименование оборудования, сооружения или строения",
		"Дата начала",
		"Срок по плану",
		"Ф.И.О механика",
		"Номер деф.вед"},
	"defectlist": {"Рег. №",
		"Наименование оборудования, сооружения или строения",
		"Описание дефектов с указанием единицы измерения и объема работ",
		"Время начала ремонта",
		"Время окончания ремонта",
		"Продолжительность ремонта (в днях)",
	},
	"report": {"Вид ремонта",
		"Количество выполненных ремонтов",
		"Кол-во дней по плану",
		"Кол-во дней факт",
		"Количество ремонтов, выполненных с нарушением сроков",
	},
	"check": {"Дата установки",
		"Установлено счетчиков",
		"Сумма оплаты",
		"Выполнена поверка",
		"Сумма оплаты",
	},
	"meters": {"номер лицевого счёта",
		"номер счётчика",
		"ФИО",
		"Адрес",
		"Номер телефона",
	},
	"checks": {"номер лицевого счёта",
		"ФИО",
		"Адрес",
		"номер счётчика",
		"Дата поверки",
		"Сумма оплаты, руб",
	},
	"installed": {"Номер лицевого счёта",
		"ФИО",
		"Адрес",
		"Номер счётчика",
		"Номер телефона",
		"Дата установки",
		"Сумма оплаты, руб",
	},
	"will_check": {"Номер лицевого счёта",
		"ФИО",
		"Адрес",
		"Номер счётчика",
		"Номер телефона",
		"Дата поверки",
	},
}

func (s *Server) getReportData(option string, form []string) (title []string, data [][]string) {

	title = TITLES[option]
	query := ""
	if len(form) == 1 {
		query = fmt.Sprintf(QUERIES[option], form[0])
	} else if len(form) == 2 {
		query = fmt.Sprintf(QUERIES[option], form[0], form[1])
	}

	rows, err := s.db.Queryx(query)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		slice, err := rows.SliceScan()
		var strSlice []string

		if err != nil {
			logrus.Error(err)
			return
		}

		for i := range slice {
			switch v := slice[i].(type) {
			case string:
				strSlice = append(strSlice, v)
			case int:
				strSlice = append(strSlice, fmt.Sprint(v))
			default:
				strSlice = append(strSlice, fmt.Sprint(v))
			}
		}

		data = append(data, strSlice)
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return
	}

	return
}
