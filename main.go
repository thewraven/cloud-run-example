package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		log.Fatalln("Es requerido especificar ENV(PORT) donde escuchará el servidor http")
	}
	addr := ":" + port
	http.HandleFunc("/", comprimirImagen)
	log.Println("Iniciando servidor http en", addr)
	err := http.ListenAndServe(addr, nil)
	log.Println(err)
}

func comprimirImagen(w http.ResponseWriter, r *http.Request) {
	//decodifica la imagen recibida en la peticion
	image, mime, err := image.Decode(r.Body)
	if err != nil {
		msg := kv{
			"error": fmt.Sprintf("No se logro decodificar la peticion como imagen: %s", err.Error()),
		}
		escribirJSON(msg, http.StatusBadRequest, w)
		return
	}
	if mime != "jpeg" {
		msg := kv{
			"error": fmt.Sprintf("Solo se soportan imagenes JPEG. formato recibido: %s", mime),
		}
		escribirJSON(msg, http.StatusBadRequest, w)
		return
	}
	//obtiene el parametro de la calidad a usar
	c := r.URL.Query().Get("calidad")
	calidad := 75
	if strings.TrimSpace(c) != "" {
		calidad, err = strconv.Atoi(c)
		if err != nil {
			msg := kv{
				"error": fmt.Sprintf("El parametro calidad solo puede ser numerico: %s", err.Error()),
			}
			escribirJSON(msg, http.StatusBadRequest, w)
			return
		}
	}
	//comprime la imagen y la escribe en un buffer
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image, &jpeg.Options{
		Quality: calidad,
	})
	if err != nil {
		msg := kv{
			"error": fmt.Sprintf("Error al codificar imagen: %s", err.Error()),
		}
		escribirJSON(msg, http.StatusInternalServerError, w)
		return
	}
	//si todo sale bien responde con la imagen comprimida
	w.Header().Set("Content-Type", "image/"+mime)
	w.WriteHeader(http.StatusOK)
	respuesta, _ := io.Copy(w, buf)
	log.Println("Tamaño de imagen entregada:", respuesta)
}

func escribirJSON(valores kv, codigo int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	err := json.NewEncoder(w).Encode(valores)
	return err
}

type kv map[string]interface{}
