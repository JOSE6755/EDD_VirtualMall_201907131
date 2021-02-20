package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
)

type Datos_fin struct {
	Datos []Datoss `json:"Datos"`
}

type Datoss struct {
	Indice        string          `json:"Indice"`
	Departamentos []Departamentos `json:"Departamentos"`
}

type Departamentos struct {
	Nombre  string    `json:"Nombre"`
	Tiendas []Tiendas `json:"Tiendas"`
}

type Tiendas struct {
	Nombre       string `json:"Nombre"`
	Descripcion  string `json:"Descripcion"`
	Contacto     string `json:"Contacto"`
	Calificacion int    `json:"Calificacion"`
}

type busqueda struct {
	Departamento string `json:"Departamento"`
	Nombre       string `json:"Nombre"`
	Calificacion int    `json:"Calificacion"`
}

type eliminacion struct {
	Nombre       string `json:"Nombre"`
	Categoria    string `json:"Categoria"`
	Calificacion int    `json:"Calificacion"`
}

var indices []string
var depas []string
var vector []Lista_doble

func inicio(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hola")

}

func crear(w http.ResponseWriter, r *http.Request) {
	var datos Datos_fin
	reqbody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, "Inserte datos validos")

	}
	json.Unmarshal(reqbody, &datos)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(datos)

	llenar(datos)
	graficar()

}

func buscar(w http.ResponseWriter, r *http.Request) {
	var dat busqueda
	reqbody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, "Inserte datos validos")
	}
	json.Unmarshal(reqbody, &dat)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusCreated)
	encontrado(dat, w)
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", inicio)

	router.HandleFunc("/meter", crear).Methods("POST")
	router.HandleFunc("/TiendaEspecifica", buscar).Methods("POST")
	router.HandleFunc("/Eliminar", eliminar)
	router.HandleFunc("/id/{id}", pornum).Methods("GET")
	router.HandleFunc("/Eliminar", eliminar).Methods("POST")
	router.HandleFunc("/Guardar", guardar).Methods("GET")
	http.ListenAndServe(":3000", router)
	log.Fatal(http.ListenAndServe(":3000", router))

}

type nodo struct {
	/*nombre       string `json:"Nombre"`
	descripcion  string `json:"Descripcion"`
	contacto     string `json:"Contacto"`
	calificacion int    `json:"Calificacion"`
	*/
	Tiendas   Tiendas
	siguiente *nodo
	anterior  *nodo
}

type Lista_doble struct {
	inicio   *nodo
	fin      *nodo
	cantidad int
}

func (l *Lista_doble) insertar(n Tiendas) {
	nuevo := &nodo{Tiendas: n}
	if l.inicio == nil {

		l.inicio = nuevo
		l.fin = nuevo
		l.cantidad++
	} else {
		fin := l.fin
		fin.siguiente = nuevo
		fin.siguiente.anterior = fin
		l.fin = nuevo
		l.cantidad++
	}

}
func (l *Lista_doble) listar(w http.ResponseWriter) {
	inicio := l.inicio

	for inicio != nil {
		json.NewEncoder(w).Encode(inicio.Tiendas)
		inicio = inicio.siguiente
	}
}

func find(a busqueda, c Lista_doble, n int, w http.ResponseWriter) bool {
	encontrado := false
	inicio := c.inicio

	for inicio != nil {
		if inicio.Tiendas.Nombre == a.Nombre {
			encontrado = true
			fmt.Println("Si entre")
			json.NewEncoder(w).Encode(inicio.Tiendas)
			break
		} else {
			inicio = inicio.siguiente
		}
	}
	return encontrado
}

func pornum(w http.ResponseWriter, r *http.Request) {
	ind := mux.Vars(r)
	id, err := strconv.Atoi(ind["id"])

	if err != nil {
		fmt.Fprint(w, "Ingrese un dato valido")
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusCreated)
	vector[id-1].listar(w)
}

func indi(a Datos_fin) int {
	var n = 0
	indices = make([]string, len(a.Datos))
	for i := 0; i < len(indices); i++ {
		indices[i] = a.Datos[i].Indice
	}

	return n
}

func llenar(a Datos_fin) {

	depas = make([]string, len(a.Datos[0].Departamentos))
	for i := 0; i < len(depas); i++ {
		depas[i] = a.Datos[0].Departamentos[i].Nombre
	}
	indi(a)
	fmt.Println(depas)

	vector = make([]Lista_doble, (len(a.Datos) * len(a.Datos[0].Departamentos) * 5))

	ingresar(a)

}

func encontrado(a busqueda, w http.ResponseWriter) {
	n := 0
	d := 0

	for i := 0; i < len(depas); i++ {
		fmt.Println(a.Nombre)
		if depas[i] == a.Departamento {
			d = i

			break
		}
	}
	fmt.Println(len(indices))
	for i := 0; i < len(indices); i++ {

		n = (i*len(depas)+d)*5 + (a.Calificacion - 1)
		fmt.Println(n)
		if find(a, vector[n], n, w) == true {
			fmt.Println("asdasdasd")
			break
		}
	}

}

func ingresar(datos Datos_fin) {
	for i := 0; i < len(datos.Datos); i++ {
		for j := 0; j < len(datos.Datos[i].Departamentos); j++ {
			for k := 0; k < len(datos.Datos[i].Departamentos[j].Tiendas); k++ {

				vector[((i*len(datos.Datos[i].Departamentos)+j)*5 + (datos.Datos[i].Departamentos[j].Tiendas[k].Calificacion - 1))].insertar(datos.Datos[i].Departamentos[j].Tiendas[k])
			}
		}

	}
	for i := 0; i < len(vector); i++ {
		fmt.Println(i)
		n := vector[i].inicio
		for n != nil {
			fmt.Println(n.Tiendas)
			n = n.siguiente
		}
	}

}

func eliminar(w http.ResponseWriter, r *http.Request) {
	var dat eliminacion
	d := 0
	n := 0
	reqbody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, "Inserte datos validos")
	}
	json.Unmarshal(reqbody, &dat)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusCreated)
	for i := 0; i < len(depas); i++ {
		if depas[i] == dat.Categoria {
			d = i
			break
		}
	}
	for i := 0; i < len(indices); i++ {
		n = (i*len(depas)+d)*5 + (dat.Calificacion - 1)
		if vector[n].elimencon(dat, w) == true {
			break
		}
	}
	z := vector[n].inicio
	for z != nil {
		fmt.Println(z.Tiendas)
		z = z.siguiente

	}
}

func (l *Lista_doble) elimencon(a eliminacion, w http.ResponseWriter) bool {
	inicio := l.inicio
	echo := false
	if inicio != nil {
		if inicio.Tiendas.Nombre == a.Nombre && inicio.Tiendas.Calificacion == a.Calificacion {

			l.inicio = inicio.siguiente
			json.NewEncoder(w).Encode(l.inicio.Tiendas)
			echo = true
			l.cantidad--
		} else {
			for inicio != nil {
				if inicio.Tiendas.Nombre == a.Nombre && inicio.Tiendas.Calificacion == a.Calificacion {
					inicio.anterior.siguiente = inicio.siguiente
					if inicio.siguiente != nil {
						inicio.siguiente.anterior = inicio.anterior
					}
				}
				inicio = inicio.siguiente

			}
			l.cantidad--
			echo = true
		}
	}

	return echo
}

func graficar() {
	n := 0
	ayu := 0
	ayu2 := 0
	nodos := "{rank=same;"
	lista := ""
	doc := "digraph G {\nnode[shape=record]\n" + `graph[splines="ortho"]` + "\n"
	aux := make([]string, (len(indices) * len(depas) * 5))
	aux2 := make([]string, len(aux)/5)
	for i := 0; i < len(indices); i++ {
		for j := 0; j < len(depas); j++ {
			for k := 0; k < 5; k++ {
				n = (i*len(depas)+j)*5 + k

				aux[n] = "nodo" + strconv.Itoa(n) + `[label="` + indices[i] + "|" + depas[j] + "|" + "POS:" + strconv.Itoa(n+1) + `"]`

			}
		}

	}

	for i := 0; i < len(aux); i++ {
		aux2[ayu] += aux[i] + "\n"
		if ayu2 <= 4 {
			nodos += "nodo" + strconv.Itoa(i) + ";"
		}

		ayu2++
		if ayu2 == 5 {
			nodos += "}"
			aux2[ayu] += nodos
			nodos = "{rank=same;"
			ayu2 = 0
			ayu++
		}
	}
	ayu = 0
	ayu2 = 0
	nodos = ""
	dots := 0
	for i := 0; i < len(aux); i++ {
		if i != len(aux) && ayu2 < 4 {
			nodos += "nodo" + strconv.Itoa(i) + "->nodo" + strconv.Itoa(i+1) + "\n"
		}
		lista += vector[i].listar2(i)
		ayu2++

		if ayu2 == 5 {

			ayu2 = 0
			doc += aux2[ayu] + "\n" + nodos
			doc += lista + "\n}"

			fmt.Println(doc)
			err := ioutil.WriteFile("Tiendas"+strconv.Itoa(dots+1)+".dot", []byte(doc), 0644)
			if err != nil {
				log.Fatal(err)
			}
			ruta, _ := exec.LookPath("dot")
			cmd, _ := exec.Command(ruta, "-Tpng", "./Tiendas"+strconv.Itoa(dots+1)+".dot").Output()
			mode := int(0777)
			ioutil.WriteFile("Tiendas"+strconv.Itoa(dots+1)+".png", cmd, os.FileMode(mode))
			doc = "digraph G {\nnode[shape=record]\n" + `graph[splines="ortho"]` + "\n"
			nodos = ""
			lista = ""
			ayu++
			dots++
		}
	}

}

func (l *Lista_doble) listar2(n int) string {
	inicio := l.inicio
	nodos := "nodo" + strconv.Itoa(n) + "->"
	datos := ""
	if inicio != nil {
		for inicio != nil {
			datos += inicio.Tiendas.Nombre + `[label="` + inicio.Tiendas.Nombre + "|" + inicio.Tiendas.Contacto + "|" + strconv.Itoa(inicio.Tiendas.Calificacion) + `"]` + "\n"
			if inicio.siguiente != nil {
				nodos += inicio.Tiendas.Nombre + "->" + inicio.siguiente.Tiendas.Nombre
			}
			if l.inicio == l.fin {
				nodos += inicio.Tiendas.Nombre
			}
			inicio = inicio.siguiente
		}
		datos += nodos + "\n"
		return datos
	}
	return datos
}

func guardar(w http.ResponseWriter, r *http.Request) {

	var datos2 Datos_fin
	tam := 0
	datos2.Datos = make([]Datoss, len(indices))

	for i := 0; i < len(indices); i++ {
		datos2.Datos[i].Departamentos = make([]Departamentos, len(depas))
	}
	for i := 0; i < len(indices); i++ {
		datos2.Datos[i].Indice = indices[i]
		for j := 0; j < len(depas); j++ {
			datos2.Datos[i].Departamentos[j].Nombre = depas[j]
			for k := 0; k < 5; k++ {
				n := (i*len(depas)+j)*5 + k
				tam += vector[n].cantidad
				if k == 4 {
					datos2.Datos[i].Departamentos[j].Tiendas = make([]Tiendas, tam)
					tam = 0
					inicio := vector[n].inicio
					for z := 0; z < len(datos2.Datos[i].Departamentos[j].Tiendas); z++ {
						datos2.Datos[i].Departamentos[j].Tiendas[z].Nombre = inicio.Tiendas.Nombre
						datos2.Datos[i].Departamentos[j].Tiendas[z].Descripcion = inicio.Tiendas.Descripcion
						datos2.Datos[i].Departamentos[j].Tiendas[z].Contacto = inicio.Tiendas.Contacto
						datos2.Datos[i].Departamentos[j].Tiendas[z].Calificacion = inicio.Tiendas.Calificacion
						if inicio.siguiente != nil {
							inicio = inicio.siguiente
						} else {
							break
						}
					}
				}

			}

		}
	}

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(datos2)
	file, err := os.Create("datos2.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	io.Copy(file, buf)

}
