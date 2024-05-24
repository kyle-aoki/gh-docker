package main

import "net/http"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, tasty"))
	})

	check(http.ListenAndServe("0.0.0.0:8080", http.DefaultServeMux))

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
