package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Kripik pedas", Harga: 1500, Stok: 3},
	{ID: 2, Nama: "Puding", Harga: 6000, Stok: 12},
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid produk id ", http.StatusBadRequest)
	}
	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Produk tidak ketemu", http.StatusNotFound)

}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	//get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	//ganti ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk id", http.StatusBadRequest)
	}

	//get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
	}

	// loop produk cari id,  ganti sesuai data dari request
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}
	http.Error(w, "Produk tidak ketemu", http.StatusNotFound)

}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk id", http.StatusBadRequest)
	}
	for i := range produk {
		if produk[i].ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Produk berhasil dihapus",
			})
			return
		}
	}
	http.Error(w, "Produk tidak ketemu", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusCreated) //201
			json.NewEncoder(w).Encode(produkBaru)

		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OKE",
			"message": "API Running",
		})
	})
	fmt.Println("Server running di localhost :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gagal Running Server")
	}

}
