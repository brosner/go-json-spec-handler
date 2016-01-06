package jsh

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDocument(t *testing.T) {

	Convey("Document Tests", t, func() {

		doc := New()

		Convey("->New()", func() {
			So(doc.JSONAPI.Version, ShouldEqual, JSONAPIVersion)
		})

		Convey("->HasErrors()", func() {
			err := &Error{Status: 400}
			addErr := doc.AddError(err)
			So(addErr, ShouldBeNil)

			So(doc.HasErrors(), ShouldBeTrue)
		})

		Convey("->HasData()", func() {
			obj, err := NewObject("1", "user", nil)
			So(err, ShouldBeNil)

			doc.Data = append(doc.Data, obj)
			So(doc.HasData(), ShouldBeTrue)
		})

		Convey("->AddObject()", func() {
			obj, err := NewObject("1", "user", nil)
			So(err, ShouldBeNil)

			err = doc.AddObject(obj)
			So(err, ShouldBeNil)
			So(len(doc.Data), ShouldEqual, 1)
		})

		Convey("->AddError()", func() {
			testError := &Error{Status: 400}

			Convey("should successfully add a valid error", func() {
				err := doc.AddError(testError)
				So(err, ShouldBeNil)
				So(len(doc.Errors), ShouldEqual, 1)
			})

			Convey("should error if validation fails while adding an error", func() {
				badError := &Error{
					Title:  "Invalid",
					Detail: "So badly",
				}

				err := doc.AddError(badError)
				So(err, ShouldNotBeNil)
				So(doc.Errors, ShouldBeEmpty)
			})
		})

		Convey("->Build()", func() {

			testObject := &Object{
				ID:     "1",
				Type:   "Test",
				Status: http.StatusAccepted,
			}

			Convey("should accept an object", func() {
				doc := Build(testObject)

				So(doc.Data, ShouldResemble, List{testObject})
				So(doc.Status, ShouldEqual, http.StatusAccepted)
			})

			Convey("should accept a list", func() {
				list := List{testObject}
				doc := Build(list)

				So(doc.Data, ShouldResemble, list)
				So(doc.Status, ShouldEqual, http.StatusOK)
			})

			Convey("should accept an error", func() {
				err := &Error{Status: 500}
				doc := Build(err)

				So(doc.Errors, ShouldNotBeEmpty)
				So(doc.Status, ShouldEqual, err.Status)
			})
		})

	})

}