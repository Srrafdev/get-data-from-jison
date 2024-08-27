package box

import (
	box "box/tracker"
	"fmt"
	"net/http"
	"text/template"
)

// hundle err in home page
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found 404", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed 405", http.StatusMethodNotAllowed)
		return
	}

	tmp, err := template.ParseFiles("./website/pages/index.html")
	if err != nil {
		http.Error(w, "Internal Server 500", http.StatusInternalServerError)
		return
	}

	var dataArtist []box.Data_Execute
	var api box.Api
	box.SaveURL(&api)
	// get all data
	var tempData box.TempStruct
	box.Decode(&tempData, api.Locations)
	//box.Decode(&tempData, api.Dates)
	//box.Decode(&tempData, api.Relation)

	dataArtist = tempData.Index
	box.Decode(&dataArtist, api.Artists)
	locaAll := box.LenData(dataArtist)
	////
	if rr := r.ParseForm(); rr != nil{
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
// filter locations	
	loca := r.FormValue("loca")
	if loca != "" && loca != "all" {
		DT := box.FilterByLocaton(dataArtist, loca)
		dataArtist = DT
	}
// filter first album
	minFL := r.FormValue("mindate")
	maxFL := r.FormValue("maxdeta")
	if minFL != "" && maxFL != "" {
		DTD := box.FilterByFirstAlbum(dataArtist, minFL, maxFL)
		dataArtist = DTD
	}
// filter Creation Date
	minCD := r.FormValue("minDaCre")
	maxCD := r.FormValue("maxDaCre")
	if minCD != "" && maxCD != ""{
		DTCD := box.FilterByCreationYear(dataArtist, minCD,maxCD)
		dataArtist = DTCD
	}

	exec := box.Execute{Loca: locaAll, DataEX: dataArtist}

	if err := tmp.Execute(w, exec); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server 500", http.StatusInternalServerError)
		return
	}
}
