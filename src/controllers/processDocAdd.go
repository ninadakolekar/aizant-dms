package controllers

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"

	"github.com/ninadakolekar/go-dms/src/auth"
	"github.com/ninadakolekar/go-dms/src/constants"

	doc "github.com/ninadakolekar/go-dms/src/docs"
	model "github.com/ninadakolekar/go-dms/src/models"
	user "github.com/ninadakolekar/go-dms/src/user"
	utility "github.com/ninadakolekar/go-dms/src/utility"
)

// ProcessDocAdd ... Process the form-values and add the document
func ProcessDocAdd(w http.ResponseWriter, r *http.Request) {
	datab := false
	datamsg := "hi"
	errb := false
	if r.Method == "POST" {

		// User Auth Verification

		username, err := auth.GetCurrentUser(r)

		if err != nil { // Auth unsucessful
			fmt.Println("ERROR ProcessDocAdd Line 23: ", err) // Debug
			http.Redirect(w, r, "/", 302)
			return
		}

		user, err := user.FetchUserByUsername(username)

		if err != nil { // User fetch unsucessful
			fmt.Println("ERROR ProcessDocAdd Line 31: ", err) // Debug
			http.Redirect(w, r, "/", 302)
			return
		}

		if user.AvailableInit == false {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		initTime := utility.XMLTimeNow()

		//TODO : Sanitize the form data

		r.ParseForm()

		docNumber := html.EscapeString(r.Form["docNumber"][0])
		docName := html.EscapeString(r.Form["docName"][0])
		docProcess := html.EscapeString(r.Form["docProcess"][0])
		docType := html.EscapeString(r.Form["docType"][0])
		docDept := html.EscapeString(r.Form["docDept"][0])
		docEffDate := html.EscapeString(r.Form["docEffDate"][0])
		docExpDate := html.EscapeString(r.Form["docExpDate"][0])
		docCreator := html.EscapeString(r.Form["docCreator"][0])
		docAuth := r.Form["docAuth"]
		docReviewers := r.Form["docReviewers"]
		docApprovers := r.Form["docApprovers"]
		for i, e := range docAuth {
			docAuth[i] = html.EscapeString(e)
		}
		for i, e := range docReviewers {
			docReviewers[i] = html.EscapeString(e)
		}
		for i, e := range docApprovers {
			docApprovers[i] = html.EscapeString(e)
		}

		fmt.Println("Form Received\n ", docNumber, docName, docProcess, docType, docDept, utility.XMLDate(docEffDate), utility.XMLDate(docExpDate), docCreator, docReviewers, docApprovers, docAuth, initTime) // Debug

		// Server-side validation

		if !doc.ValidateDocNo(docNumber) && doc.ValidateDocName(docName) {

			// Make a new inactiveDoc struct using received form data

			newDoc := model.InactiveDoc{
				DocNo:           docNumber,
				Title:           docName,
				DocType:         docType,
				DocProcess:      docProcess,
				DocEffDate:      utility.XMLDate(docEffDate),
				DocExpDate:      utility.XMLDate(docExpDate),
				DocStatus:       false,
				Initiator:       username,
				Creator:         docCreator,
				Reviewer:        docReviewers,
				Approver:        docApprovers,
				Authorizer:      docAuth,
				DocDept:         docDept,
				FlowStatus:      constants.QaFlow,
				FlowList:        nil,
				CurrentFlowUser: 0,
				DocTemplate:     0,
				InitTS:          initTime,
				CreateTS:        "",
				ReviewTS:        "",
				AuthTS:          "",
				ApproveTS:       "",
				DocumentBody:    nil,
				QA:              "firefox",
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
	// http.Redirect(w, r, "/doc/add", http.StatusSeeOther)

	tmpl := template.Must(template.ParseFiles("templates/addNewDoc.html"))

	tmpl.Execute(w, docAddMsg{Datab: datab, Errb: errb, Datamsg: datamsg, Approvers: SendApprovers(), Reviewers: SendReviewers(), Authorisers: SendAuthoriser(), Creators: SendCreators(), DocumentExist: false, Redirect: true, Document: model.InactiveDoc{}})

}
