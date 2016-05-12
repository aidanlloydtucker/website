package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yosssi/ace"
)

type Class struct {
	Id   string
	Name string
	Num  int
}

type Classes []Class

func (slice Classes) Len() int {
	return len(slice)
}

func (slice Classes) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice Classes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Assignment struct {
	Name     string
	Desc     string
	Date     string
	Category string
}

var MASTER_CLASS_LIST = Classes{
	Class{Id: "eeng1", Name: "English I/II", Num: 873},
	Class{Id: "eeng3", Name: "English III/IV", Num: 624},
	Class{Id: "eengf", Name: "English Fundementals", Num: 347},
	Class{Id: "hac", Name: "Ancient Cultures", Num: 626},
	Class{Id: "hbh", Name: "British History", Num: 776},
	Class{Id: "hus", Name: "US History", Num: 565},
	Class{Id: "hws", Name: "World Studies II", Num: 564},
	Class{Id: "lasl1", Name: "ASL I", Num: 389},
	Class{Id: "lasl2", Name: "ASL II", Num: 390},
	Class{Id: "lasl3", Name: "ASL III/IV", Num: 391},
	Class{Id: "lspa1", Name: "Spanish I", Num: 416},
	Class{Id: "lspa2", Name: "Spanish II", Num: 418},
	Class{Id: "lspa3", Name: "Spanish III", Num: 420},
	Class{Id: "lspa4", Name: "Spanish IV", Num: 422},
	Class{Id: "malg1", Name: "Algebra 1", Num: 892},
	Class{Id: "malg2", Name: "Algebra II", Num: 628},
	Class{Id: "mcalc", Name: "Calculus", Num: 496},
	Class{Id: "mcons", Name: "Consumer Math", Num: 149},
	Class{Id: "mgeo", Name: "Geometry", Num: 955}, //was 893
	Class{Id: "mpalg", Name: "Introduction to Algebra", Num: 500},
	Class{Id: "mtrig", Name: "Trigonometry", Num: 894},
	Class{Id: "oart1", Name: "Studio Art", Num: 902},
	Class{Id: "oart2", Name: "Art II", Num: 901},
	Class{Id: "odra", Name: "Drama", Num: 930},
	Class{Id: "ogov", Name: "Government", Num: 342},
	Class{Id: "ogs", Name: "Gender Studies", Num: 715},
	Class{Id: "ohr", Name: "Human Relations", Num: 783},
	Class{Id: "omsic", Name: "Music", Num: 623},
	Class{Id: "osocs", Name: "Socratic Seminar", Num: 457},
	Class{Id: "otsd", Name: "Transitions Seminar (Deborah C.)", Num: 903},
	Class{Id: "otsj", Name: "Transitions Seminar (Janel C.)", Num: 351},
	Class{Id: "sabio", Name: "Advanced Biology", Num: 427},
	Class{Id: "sbio", Name: "Biology", Num: 430},
	Class{Id: "schem", Name: "Chemistry", Num: 960},
	Class{Id: "sphys", Name: "Physics", Num: 888},
	Class{Id: "ssci1", Name: "Science 1", Num: 886},
	Class{Id: "none", Name: "Other/None", Num: 0},
}

func HomeworkHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	cl, ok := session.Values["class_list"]
	if !ok {
		session.Save(r, w)
		http.Redirect(w, r, "/homework/classes", http.StatusFound)
		return
	}

	tpl, err := p.Load("base", "homework", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reflect.TypeOf(cl).Kind() == reflect.Array || reflect.TypeOf(cl).Kind() == reflect.Slice {

		valuedCl := reflect.ValueOf(cl)
		var actClassList Classes

		for i := 0; i < valuedCl.Len(); i++ {
			for _, cc := range MASTER_CLASS_LIST {
				if valuedCl.Index(i).Int() == int64(cc.Num) {
					actClassList = append(actClassList, cc)
				}
			}
		}

		data := map[string]interface{}{
			"Classes": actClassList,
		}

		if err := tpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.Redirect(w, r, "/homework/classes", http.StatusFound)
		return
	}

}

func HomeworkAssignmentsHandler(w http.ResponseWriter, r *http.Request) {
	classnum := r.URL.Query().Get("classnum")

	doc, err := goquery.NewDocument("http://www.mid-pen.com/Page/" + classnum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	assignArr := doc.Find("h2").Filter(":contains('Current Assignments')").NextFiltered(".ui-articles").Find(".ui-article-title a")
	assignJSON := make([]Assignment, 0)
	assignArr.Each(func(i int, s *goquery.Selection) {
		category := s.Parent().Next().Next().Text()
		onclick, ok := s.Attr("onclick")
		if !ok {
			return
		}
		idNumArray := strings.Split(strings.Replace(onclick[18:len(onclick)-2], "(", ",", -1), ",")

		infoDoc, err2 := goquery.NewDocument("http://www.mid-pen.com//site/UserControls/Assignment/AssignmentViewWrapper.aspx?AssignmentID=" + idNumArray[1] + "&ModuleInstanceID=" + idNumArray[0] + "&PageID=" + idNumArray[2])
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}

		desc, err2 := infoDoc.Find(".ui-article-body").Html()
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}

		date := infoDoc.Find(".ui-article-detail").Last().Text()

		if strings.Contains(desc, "/cms/") {
			desc = strings.Replace(desc, "/cms/", "http://mid-pen.com/cms/", -1)
		}
		assignJSON = append(assignJSON, Assignment{Name: s.Text(), Desc: desc, Date: date, Category: category})
	})

	jsn, err := json.Marshal(assignJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)

}

// TODO: Make not sort every load
func HomeworkGETClassesHandler(w http.ResponseWriter, r *http.Request) {

	tpl, err := p.Load("views/base", "views/classes", &ace.Options{
		FuncMap: template.FuncMap{
			"SelectIf": func(prevClasses []string, classPer int, id string) bool {
				return len(prevClasses) > classPer && id == prevClasses[classPer]
			},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, sessionName)
	cl, ok := session.Values["class_list"]

	sortedClassList := make(Classes, len(MASTER_CLASS_LIST))
	copy(sortedClassList, MASTER_CLASS_LIST)

	sort.Sort(sortedClassList)

	data := map[string]interface{}{
		"Classes": sortedClassList,
	}

	if ok && (reflect.TypeOf(cl).Kind() == reflect.Array || reflect.TypeOf(cl).Kind() == reflect.Slice) {

		valuedCl := reflect.ValueOf(cl)
		var prevClasses []string

		for i := 0; i < valuedCl.Len(); i++ {
			for _, cc := range MASTER_CLASS_LIST {
				if valuedCl.Index(i).Int() == int64(cc.Num) {
					prevClasses = append(prevClasses, cc.Id)
				}
			}
		}

		data["IsClasses"] = true
		data["PrevClasses"] = prevClasses
	} else {
		data["IsClasses"] = false
		data["PrevClasses"] = []string{}
	}

	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HomeworkPUTClassesHandler(w http.ResponseWriter, r *http.Request) {
	var classlist []string
	json.Unmarshal([]byte(r.FormValue("classlist")), &classlist)

	var classesByNum []int64

	for _, cl := range classlist {
		for _, cc := range MASTER_CLASS_LIST {
			if cl == cc.Id {
				classesByNum = append(classesByNum, int64(cc.Num))
			}
		}
	}

	session, _ := store.Get(r, sessionName)
	session.Values["class_list"] = classesByNum
	session.Save(r, w)
}

func HomeworkAidanHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := p.Load("base", "homework", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aidanClassesById := []string{"mgeo", "schem", "lspa1", "odra", "eeng1", "hac"}
	var aidanClasses Classes

	for _, cl := range aidanClassesById {
		for _, cc := range MASTER_CLASS_LIST {
			if cl == cc.Id {
				aidanClasses = append(aidanClasses, cc)
			}
		}
	}

	data := map[string]interface{}{
		"Classes": aidanClasses,
	}

	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
