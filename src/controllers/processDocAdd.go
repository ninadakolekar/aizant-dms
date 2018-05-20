package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	doc "github.com/ninadakolekar/aizant-dms/src/docs"
	model "github.com/ninadakolekar/aizant-dms/src/models"
	utility "github.com/ninadakolekar/aizant-dms/src/utility"
)

// ProcessDocAdd ... Process the form-values and add the document
func ProcessDocAdd(w http.ResponseWriter, r *http.Request) {
	datab := false
	datamsg := "hi"
	errb := false
	if r.Method == "POST" {

		initTime := utility.XMLTimeNow()

		//TODO : Sanitize the form data

		r.ParseForm()

		docNumber := r.Form["docNumber"][0]
		docName := r.Form["docName"][0]
		docProcess := r.Form["docProcess"][0]
		docType := r.Form["docType"][0]
		docDept := r.Form["docDept"][0]
		docEffDate := r.Form["docEffDate"][0]
		docExpDate := r.Form["docExpDate"][0]
		docCreator := r.Form["docCreator"][0]
		docAuth := r.Form["docAuth"]
		docReviewers := r.Form["docReviewers"]
		docApprovers := r.Form["docApprovers"]

		fmt.Println("Form Received\n ", docNumber, docName, docProcess, docType, docDept, utility.XMLDate(docEffDate), utility.XMLDate(docExpDate), docCreator, docReviewers, docApprovers, docAuth, initTime) // Debug

		// Server-side validation

		if !doc.ValidateDocNo(docNumber) && doc.ValidateDocName(docName) {
			// Make a new inactiveDoc struct using received form data

			newDoc := model.InactiveDoc{
				DocNo:        docNumber,
				Title:        docName,
				DocType:      docType,
				DocProcess:   docProcess,
				DocEffDate:   utility.XMLDate(docEffDate),
				DocExpDate:   utility.XMLDate(docExpDate),
				DocStatus:    false,
				Initiator:    "self", // Initiator is "self" currently
				Creator:      docCreator,
				Reviewer:     docReviewers,
				Approver:     docApprovers,
				Authorizer:   docAuth,
				DocDept:      docDept,
				FlowStatus:   0,
				DocTemplate:  0,
				InitTS:       initTime,
				CreateTS:     "",
				ReviewTS:     "",
				AuthTS:       "",
				ApproveTS:    "",
				DocumentBody: []string{"Empty Body"},
			}
			// Insert the new document
			resp, err := doc.AddInactiveDoc(newDoc)

			// Respond
			if err != nil {
				// fmt.Println("ERROR ProcessDocAdd() Line 47: " + err.Error()) // Debug
				errb = true
				datamsg = "Failed to create new document"

			} else {
				log.Println(resp) // Debug
				datab = true
				datamsg = "Successfully Intiated New Document"
			}

		} else {
			errb = true
			datamsg = "Failed to create new document (Invalid Document Number or Name)."
		}
	}

	// Render a new form
	tmpl := template.Must(template.ParseFiles("templates/addNewDoc.html"))

	tmpl.Execute(w, templateData{Datab: datab, Errb: errb, Datamsg: datamsg, Approvers: SendApprovers(), Reviewers: SendReviewers(), Authorisers: SendAuthoriser(), Creators: SendCreators()})
}
