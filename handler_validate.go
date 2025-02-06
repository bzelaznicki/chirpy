package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)

		return
	}

	valid, msg := validateLength(params.Body)

	if !valid {
		respondWithError(w, http.StatusBadRequest, msg, nil)
		return
	}

	cleanedChirp, err := profanityChecker(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to proces chirp.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedChirp, //`json:"cleaned_body"`
	})

}

func validateLength(s string) (bool, string) {
	if len(s) > maxChirpLength {
		return false, "Chirp too long"
	}
	if len(s) == 0 {
		return false, "Chirp empty"
	}
	return true, ""
}

func profanityChecker(s string) (string, error) {

	if len(s) == 0 {
		return "", fmt.Errorf("chirp is empty")
	}
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitString := strings.Split(s, " ")

	for i, word := range splitString {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				splitString[i] = "****"
			}
		}
	}

	mergedString := strings.Join(splitString, " ")

	return mergedString, nil
}
