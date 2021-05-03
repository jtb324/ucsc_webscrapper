package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func format_query(gene_id string) []byte {
	/*function to format the graphQL query
	Parameters
	__________
	gene_id string
		string containing the id of the gene

	Returns
	_______
	[]byte
		returns a slice of bytes where the query dictionary was converted
		to a byte object
	*/
	jsonData := map[string]string{
		"query": `
			{
				gene(gene_symbol: "` + gene_id + `", reference_genome: GRCh37) {
		 			start
		 			stop
					omim_id
    				name
    				chrom
				}
		  	}
		  `,
	}

	jsonValue, _ := json.Marshal(jsonData)

	return jsonValue
}

type Gene struct {
	Gene gene_info `json:"gene"`
}
type gene_info struct {
	Start   int    `json:"start"`
	Stop    int    `json:"stop"`
	Omim_id string `json:"omim_id"`
	Name    string `json:"name"`
	Chrom   int    `json:"chrom"`
}

type Data struct {
	Data Gene `json:"data"`
}

func fetch_response(api_website string, gene_list []string) {
	/*function to fetch the reponse from the the bnomad api
	Parameters
	__________
	api_website string
		url to the api of interest

	gene_list []string
		slice of strings that has each gene id
	*/
	request_made := 0

	var gene_info_slice []Data

	for i := 0; i < len(gene_list); i++ {

		jsonByteString := format_query(gene_list[i])
		// fmt.Println(jsonByteString)
		request, error := http.NewRequest("POST", api_website, bytes.NewBuffer(jsonByteString))

		if error != nil {
			log.Fatalln(error)
		}
		request.Header.Add("Content-Type", "application/json")

		client := &http.Client{Timeout: time.Second * 10}

		response, response_err := client.Do(request)

		if response_err != nil {
			log.Fatalf("The HTTP request failed with error %s\n", response_err)
		}

		//deferiing the responses close
		defer response.Body.Close()

		data, _ := ioutil.ReadAll(response.Body)

		//convert the data to a json object
		var json_response Data

		json.Unmarshal(data, &json_response)

		fmt.Println(json_response.Data.Gene.Name)

		//creating a slice that has all the gene information from the api
		gene_info_slice = append(gene_info_slice, json_response)
		//updating request counter
		request_made++
		// making the program sleep for a second after every three requests
		if request_made%4 == 0 {
			time.Sleep(time.Second)
		}
	}
	fmt.Println(len(gene_info_slice))
}
