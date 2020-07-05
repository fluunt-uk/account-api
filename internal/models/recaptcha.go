package models

type ReCaptcha struct {
	Success  	bool 		`json:"success"`
	Error    	[]string 	`json:"error-codes"`
}

