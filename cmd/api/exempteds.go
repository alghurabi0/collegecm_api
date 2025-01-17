package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"collegecm.hamid.net/internal/data"
	"collegecm.hamid.net/internal/validator"
)

func (app *application) getExempteds(w http.ResponseWriter, r *http.Request) {
	exempteds, err := app.models.Exempteds.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"exempteds": exempteds}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getExempted(w http.ResponseWriter, r *http.Request) {
	//id
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	exempted, err := app.models.Exempteds.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"exempted": exempted}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) findExempted(w http.ResponseWriter, r *http.Request) {
	student_idStr := r.PathValue("student_id")
	if strings.TrimSpace(student_idStr) == "" {
		app.notFoundResponse(w, r)
		return
	}
	student_id, err := strconv.ParseInt(student_idStr, 10, 64)
	if err != nil || student_id < 1 {
		app.notFoundResponse(w, r)
		return
	}
	subject_idStr := r.PathValue("subject_id")
	if strings.TrimSpace(subject_idStr) == "" {
		app.notFoundResponse(w, r)
		return
	}
	subject_id, err := strconv.ParseInt(subject_idStr, 10, 64)
	if err != nil || subject_id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	exempted, err := app.models.Exempteds.Find(student_id, subject_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"exempted": exempted}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSubjectsExempteds(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	exempteds, err := app.models.Exempteds.GetSubjects(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"subjects_exempteds": exempteds}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getStudentsExempteds(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	exempteds, err := app.models.Exempteds.GetStudents(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"students_exempteds": exempteds}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createExempted(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StudentId int64 `json:"student_id"`
		SubjectId int64 `json:"subject_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	exempted := &data.Exempted{
		StudentId: input.StudentId,
		SubjectId: input.SubjectId,
	}
	v := validator.New()
	if data.ValidateExempted(v, exempted); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Exempteds.Insert(exempted)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	student, err := app.models.Students.Get(exempted.StudentId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	subject, err := app.models.Subjects.Get(exempted.SubjectId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	exempted.StudentName = student.StudentName
	exempted.SubjectName = subject.SubjectName
	err = app.writeJSON(w, http.StatusCreated, envelope{"exempted": exempted}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteExempted(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Exempteds.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "تم الحذف بنجاح"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
