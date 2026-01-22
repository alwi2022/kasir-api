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

type Category struct {
  ID          int    `json:"id"`
  Name        string `json:"name"`
  Description string `json:"description"`
}


var categories = []Category{
  {ID: 1, Name: "Makanan Ringan", Description: "Camilan lezat untuk menemani hari Anda"},
  {ID: 2, Name: "Minuman", Description: "Segarkan dahaga Anda dengan minuman kami"},
}

var produk = []Produk{
	{ID: 1, Nama: "Kripik pedas", Harga: 1500, Stok: 3},
	{ID: 2, Nama: "Puding", Harga: 6000, Stok: 12},
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid category id ", http.StatusBadRequest)
		return
	}
	for _, k := range categories {
		if k.ID == id {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(k)
			return
		}
	}
	http.Error(w, "Category tidak ketemu", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	//get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	//ganti ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category id", http.StatusBadRequest)
		return
	}
	//get data dari request
	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// loop category cari id,  ganti sesuai data dari request
	for i := range categories {
		if categories[i].ID == id {
			updateCategory.ID = id
			categories[i] = updateCategory
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}
	http.Error(w, "Category tidak ketemu", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category id", http.StatusBadRequest)
		return
	}
	for i := range categories {
		if categories[i].ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category berhasil dihapus",
			})
			return
		}
	}
	http.Error(w, "Category tidak ketemu", http.StatusNotFound)
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid produk id ", http.StatusBadRequest)
		return
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
		return
	}

	//get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
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
		return
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
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			var categoryBaru Category
			err := json.NewDecoder(r.Body).Decode(&categoryBaru)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}
			categoryBaru.ID = len(categories) + 1
			categories = append(categories, categoryBaru)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(categoryBaru)
		}

	})

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello BANG"))
	})

	fmt.Println("Server running di localhost :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gagal Running Server")
	}

}
