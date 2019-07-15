package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var cli *client.Client

func main() {

	var err error
	cli, err = client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getContainers(w)
	})

	if err := http.ListenAndServe(":8282", nil); err != nil {
		panic(err)
	}

	for {
		<-time.After(10 * time.Second)
		go getImages()
	}

}

func getContainers(w http.ResponseWriter) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	colTitles := []string{"Image", "Ports", "State"}
	info := make(map[string][]string)
	for _, container := range containers {
		ports := ""
		for _, port := range container.Ports {
			ports += utilConcat([]string{
				"[", port.Type, "=",
				strconv.Itoa(int(port.PrivatePort)), ":",
				strconv.Itoa(int(port.PublicPort)), "]<br>"})
		}
		info[container.ID[:10]] = []string{container.ID[:10], ports, container.State}
	}
	fmt.Fprintf(w, "<html>"+getTable("Containers List", colTitles, info)+"</html>")
}

func getImages() {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID, image.RepoTags, image.Size, image.Containers)
	}
}

func getTable(title string, colums []string, rows map[string][]string) string {

	var colTitles string = ""
	for _, col := range colums {
		colTitles += utilConcat([]string{"<th>", col, "</th>"})
	}

	var rowContent string = ""
	for _, v := range rows {
		rowContent += "<tr>"
		for _, col := range v {
			rowContent += getCel(col)
		}
		rowContent += "</tr>"
	}

	return "<h3>" + title + "</h3>" + "<table>" + "<tr>" + colTitles + "</tr>" + "<tr>" + rowContent + "</tr>" + "</table>"
}

func getCel(value string) string {
	return "<th>" + value + "</th>"
}

func utilConcat(content []string) string {
	var buffer bytes.Buffer
	for _, text := range content {
		fmt.Println(text)
		buffer.WriteString(text)
	}
	return buffer.String()
}
