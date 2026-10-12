package main

import (
	"http"
);

type Game struct {
	ID		int		`json:"_id"`
	Title		string		`json:"title"`
	RequiredAge	any		`json:"required_age"`
	IsFree
}

function main() {


