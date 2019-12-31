package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid: true,
	}
}

type GenDateVal struct {
	Year int
	Month int
	Day int
}

type GenDate struct {
	Type      string
	StartDate *GenDateVal
	EndDate   *GenDateVal
}

// Value does the conversion to a string
func (g GenDate) Value() (driver.Value, error) {

	startDate := "\"(0,0,0,0,0)\""
	if(g.StartDate != nil) {
		startDate = "\"(0,0," + strconv.Itoa(g.StartDate.Year) +
			"," + strconv.Itoa(g.StartDate.Month) +
			"," + strconv.Itoa(g.StartDate.Day) +")\""
	}
	endDate := "\"(0,0,0,0,0)\""
	if(g.EndDate != nil) {
		endDate = "\"(0,0," + strconv.Itoa(g.EndDate.Year) +
			"," + strconv.Itoa(g.EndDate.Month) +
			"," + strconv.Itoa(g.EndDate.Day) +")\""
	}
	data := "(" + g.Type + "," + startDate +"," + endDate + ",0)"
	log.Printf("Converted GenDate to a string: %s", data)
	return data, nil
}

type dateval struct {
	Text     string `xml:",chardata"`
	Modifier string `xml:"Modifier"`
	Year     string `xml:"Year"`
	Month    string `xml:"Month"`
	Day      string `xml:"Day"`
}

type date struct {
	Text      string `xml:",chardata"`
	Type      string `xml:"Type,attr"`
	StartDate struct {
		Text    string `xml:",chardata"`
		DateVal  dateval `xml:"DateVal"`
	} `xml:"StartDate"`
	EndDate struct {
		Text    string `xml:",chardata"`
		DateVal dateval `xml:"DateVal"`
	} `xml:"EndDate"`
}

type name struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"ID,attr"`
	Type        string `xml:"Type,attr"`
	IsPreferred string `xml:"IsPreferred"`
	PersonID    struct {
		Text string `xml:",chardata"`
		ID   string `xml:"ID,attr"`
	} `xml:"PersonID"`
	Given    string `xml:"Given"`
	Surname  string `xml:"Surname"`
	Familiar string `xml:"Familiar"`
	Date     date `xml:"Date"`
	Title     string `xml:"Title"`
	DisplayAs string `xml:"DisplayAs"`
}

type header struct {
	Text    string `xml:",chardata"`
	Created struct {
		Text string `xml:",chardata"`
		Date date `xml:"Date"`
		Version     string `xml:"Version"`
		DBInfo      string `xml:"DBInfo"`
		Copyright   string `xml:"Copyright"`
		PeopleCount string `xml:"PeopleCount"`
		FamilyCount string `xml:"FamilyCount"`
	} `xml:"Created"`
	Researcher struct {
		Text    string `xml:",chardata"`
		Contact struct {
			Text        string `xml:",chardata"`
			SimpleName  string `xml:"SimpleName"`
			AddressLine []struct {
				Text    string `xml:",chardata"`
				LineNum string `xml:"LineNum,attr"`
			} `xml:"AddressLine"`
			Email string `xml:"Email"`
			URL   string `xml:"URL"`
		} `xml:"Contact"`
	} `xml:"Researcher"`
}

type reference struct {
	Text string `xml:",chardata"`
	ID   string `xml:"ID,attr"`
}

type attachment struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"ID,attr"`
	ReferenceID reference `xml:"ReferenceID"`
	Filename string `xml:"Filename"`
	Fileinfo string `xml:"Fileinfo"`
	Detail   string `xml:"Detail"`
}

type attachments struct {
	Text       string `xml:",chardata"`
	Attachment []attachment `xml:"Attachment"`
}

type fact struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"ID,attr"`
	Type        string `xml:"Type,attr"`
	ReferenceID reference `xml:"ReferenceID"`
	Place string `xml:"Place"`
	Date  date `xml:"Date"`
	Detail string `xml:"Detail"`
	NoteID reference `xml:"NoteID"`
}

type child struct {
	Text     string `xml:",chardata"`
	ID       string `xml:"ID,attr"`
	PersonID reference `xml:"PersonID"`
	FamilyID reference `xml:"FamilyID"`
	Parent1Relation struct {
		Text         string `xml:",chardata"`
		Relationship struct {
			Text     string `xml:",chardata"`
			ParentID string `xml:"ParentID,attr"`
			Type     string `xml:"Type,attr"`
		} `xml:"Relationship"`
	} `xml:"Parent1Relation"`
	Parent2Relation struct {
		Text         string `xml:",chardata"`
		Relationship struct {
			Text     string `xml:",chardata"`
			ParentID string `xml:"ParentID,attr"`
			Type     string `xml:"Type,attr"`
		} `xml:"Relationship"`
	} `xml:"Parent2Relation"`
	Ordinal string `xml:"Ordinal"`
}

type source struct {
	Text           string `xml:",chardata"`
	ID             string `xml:"ID,attr"`
	Detail         string `xml:"Detail"`
	SourceTitle    string `xml:"SourceTitle"`
	SourceLocation string `xml:"SourceLocation"`
}

type note struct {
	Text   string `xml:",chardata"`
	ID     string `xml:"ID,attr"`
	Detail string `xml:"Detail"`
}

type family struct {
	Text    string `xml:",chardata"`
	ID      string `xml:"ID,attr"`
	UserID  string `xml:"UserID"`
	PrimeID reference `xml:"PrimeID"`
	PartnerID reference `xml:"PartnerID"`
	NoteID reference `xml:"NoteID"`
}

type person struct {
	Text     string `xml:",chardata"`
	ID       string `xml:"ID,attr"`
	UserID   string `xml:"UserID"`
	BirthSex string `xml:"BirthSex"`
	NoteID   reference `xml:"NoteID"`
	IsPrivate string `xml:"IsPrivate"`
	SourceID  reference `xml:"SourceID"`
}

type ScionPC struct {
	XMLName                   xml.Name `xml:"ScionPC"`
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"noNamespaceSchemaLocation,attr"`
	Header                    header `xml:"Header"`
	Names struct {
		Text string `xml:",chardata"`
		Name [] name `xml:"Name"`
	} `xml:"Names"`
	PersonalFacts struct {
		Text string `xml:",chardata"`
		Fact []fact `xml:"Fact"`
	} `xml:"PersonalFacts"`
	FamilyFacts struct {
		Text string `xml:",chardata"`
		Fact []fact `xml:"Fact"`
	} `xml:"FamilyFacts"`
	People struct {
		Text   string `xml:",chardata"`
		Person []person `xml:"Person"`
	} `xml:"People"`
	Families struct {
		Text   string `xml:",chardata"`
		Family []family `xml:"Family"`
	} `xml:"Families"`
	Children struct {
		Text  string `xml:",chardata"`
		Child []child `xml:"Child"`
	} `xml:"Children"`
	Notes struct {
		Text string `xml:",chardata"`
		Note []note `xml:"Note"`
	} `xml:"Notes"`
	Sources struct {
		Text   string `xml:",chardata"`
		Source []source `xml:"Source"`
	} `xml:"Sources"`
	PersonalAttachments attachments `xml:"PersonalAttachments"`
	FamilyAttachments   attachments `xml:"FamilyAttachments"`
}


func main() {
	log.SetFlags(log.Lshortfile)
	xmlFile, err := os.Open("Arnholm.sgx")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var scion ScionPC
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &scion)
	fmt.Println("Scion XML Read")

	db, err := sql.Open("postgres", "user=gofamily dbname=gofamily password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, note := range scion.Notes.Note {
		sqlInsertNote(note, db)
	}
	for _, source := range scion.Sources.Source {
		sqlInsertSource(source, db)
	}
	for _, person := range scion.People.Person {
		sqlInsertPerson(person, db)
	}
	for _, family := range scion.Families.Family {
		sqlInsertFamily(family, db)
	}
	for _, child := range scion.Children.Child {
		sqlInsertChild(child, db)
	}
	for _, attachment := range scion.PersonalAttachments.Attachment {
		sqlInsertAttachment(attachment, db)
	}
	for _, attachment := range scion.FamilyAttachments.Attachment {
		sqlInsertAttachment(attachment, db)
	}
	for _, name := range scion.Names.Name {
		sqlInsertName(name, db)
	}
	for _, fact := range scion.PersonalFacts.Fact {
		sqlInsertFact(fact, db)
	}
	for _, fact := range scion.FamilyFacts.Fact {
		sqlInsertFact(fact, db)
	}

	// Mjupp(scion)
}

func sqlInsertFact(fact fact, db *sql.DB) {
	fmt.Printf("Fact: id: %v Type %v ReferenceID %v Place %v Date %v Detail %v NoteID %v\n",
		fact.ID, fact.Type, fact.ReferenceID.ID, fact.Place, fact.Date, fact.Detail, fact.NoteID.ID)
	stmt, err := db.Prepare(
		"INSERT INTO fact(id, type, referenceid, date, place, detail, noteid)" +
			" VALUES ($1, $2, $3, $4, $5, $6, $7)" +
			" ON CONFLICT (id) " +
			"    DO UPDATE SET type = excluded.type," +
			"                  referenceid = excluded.referenceid," +
			"                  place = excluded.place," +
			"                  date = excluded.date," +
			"                  detail = excluded.detail," +
			"                  noteid = excluded.noteid;")
	if err != nil {
		log.Fatal(err)
	}
	date := dateToGenDate(fact.Date)
	_, err = stmt.Exec(fact.ID, NewNullString(fact.Type), NewNullString(fact.ReferenceID.ID), date,
		NewNullString(fact.Place), NewNullString(fact.Detail), NewNullString(fact.NoteID.ID))
	if err != nil {
		log.Fatal(err)
	}
}

func sqlInsertName(name name, db *sql.DB) {
	fmt.Printf("Name: id: %v Type %v prefered %v person %v Given %v Surname %v Familiar %v "+
		"Date %v Title %v Display %v\n",
		name.ID, name.Type, name.IsPreferred, name.PersonID.ID, name.Given, name.Surname, name.Familiar,
		name.Date, name.Title, name.DisplayAs)
	stmt, err := db.Prepare(
		"INSERT INTO name(id, type, isPrefered, personid, given, surname, familiar, title, displayAs, date)" +
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::GenDate)" +
			" ON CONFLICT (id) " +
			"    DO UPDATE SET type = excluded.type," +
			"                  isPrefered = excluded.isPrefered," +
			"                  personid = excluded.personid," +
			"                  given = excluded.given," +
			"                  surname = excluded.surname," +
			"                  familiar = excluded.familiar," +
			"                  title = excluded.title," +
			"                  displayAs = excluded.displayAs;")
	if err != nil {
		log.Fatal(err)
	}
	dateval := dateToGenDate(name.Date)
	_, err = stmt.Exec(name.ID, NewNullString(name.Type), name.IsPreferred == "true", NewNullString(name.PersonID.ID),
		NewNullString(name.Given), NewNullString(name.Surname), NewNullString(name.Familiar),
		NewNullString(name.Title), NewNullString(name.DisplayAs), dateval)
	if err != nil {
		log.Fatal(err)
	}
}

func dateToGenDate(date date) GenDate {
	var startDate *GenDateVal = nil
	if (len(date.StartDate.DateVal.Year) + len(date.StartDate.DateVal.Month) + len(date.StartDate.DateVal.Day)) != 0 {
		year, _ := strconv.Atoi(date.StartDate.DateVal.Year)
		month, _ := strconv.Atoi(date.StartDate.DateVal.Month)
		day, _ := strconv.Atoi(date.StartDate.DateVal.Day)
		startDate = &GenDateVal{
			Year:  year,
			Month: month,
			Day:   day,
		}
	}
	var endDate *GenDateVal = nil
	if (len(date.EndDate.DateVal.Year) + len(date.EndDate.DateVal.Month) + len(date.EndDate.DateVal.Day)) != 0 {
		year, _ := strconv.Atoi(date.EndDate.DateVal.Year)
		month, _ := strconv.Atoi(date.EndDate.DateVal.Month)
		day, _ := strconv.Atoi(date.EndDate.DateVal.Day)
		startDate = &GenDateVal{
			Year:  year,
			Month: month,
			Day:   day,
		}
	}
	dateval := GenDate{Type: date.Type, StartDate: startDate, EndDate: endDate}
	return dateval
}

func sqlInsertAttachment(attachment attachment, db *sql.DB) {
	fmt.Printf("Attachment: id: %v ReferenceID %v Filename %v"+
		" Fileinfo %v Detail %v\n",
		attachment.ID, attachment.ReferenceID.ID, attachment.Filename, attachment.Fileinfo, attachment.Detail)
	stmt, err := db.Prepare(
		"INSERT INTO attach(id, referenceid, filename, fileinfo, detail)" +
			" VALUES ($1, $2, $3, $4, $5)" +
			" ON CONFLICT (id) " +
			"    DO UPDATE SET referenceid = excluded.referenceid," +
			"                  filename = excluded.filename," +
			"                  fileinfo = excluded.fileinfo," +
			"                  detail = excluded.detail;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(attachment.ID, NewNullString(attachment.ReferenceID.ID), NewNullString(attachment.Filename),
		NewNullString(attachment.Fileinfo), NewNullString(attachment.Detail))
	if err != nil {
		log.Fatal(err)
	}
}

func sqlInsertChild(child child, db *sql.DB) {
	fmt.Printf("Child: id: %v personid %v familyid %v"+
		" parent1id %v parent1rel %v parent2id %v parent2rel %v ordinal %v\n",
		child.ID, child.PersonID.ID, child.FamilyID.ID,
		child.Parent1Relation.Relationship.ParentID, child.Parent1Relation.Relationship.Type,
		child.Parent2Relation.Relationship.ParentID, child.Parent2Relation.Relationship.Type,
		child.Ordinal)
	stmt, err := db.Prepare(
		"INSERT INTO child(id, personid, familyid," +
			"                     parent1id, parent1rel, parent2id, parent2rel, ordinal)" +
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8)" +
			" ON CONFLICT (id) " +
			"    DO UPDATE SET personid = excluded.personid," +
			"                  familyid = excluded.familyid," +
			"                  parent1id = excluded.parent1id," +
			"                  parent1rel = excluded.parent1rel," +
			"                  parent2id = excluded.parent2id," +
			"                  parent2rel = excluded.parent2rel," +
			"                  ordinal = excluded.ordinal;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(child.ID, NewNullString(child.PersonID.ID), NewNullString(child.FamilyID.ID),
		NewNullString(child.Parent1Relation.Relationship.ParentID), child.Parent1Relation.Relationship.Type,
		NewNullString(child.Parent1Relation.Relationship.ParentID), child.Parent1Relation.Relationship.Type,
		child.Ordinal)
	if err != nil {
		log.Fatal(err)
	}
}

func sqlInsertSource(source source, db *sql.DB) {
	fmt.Printf("Source: ID %v Detail %v\n", source.ID, source.Detail)
	stmt, err := db.Prepare(
		"INSERT INTO source(id, detail, title, location)" +
			" VALUES ($1, $2, $3, $4)" +
			" ON CONFLICT (id) " +
			"    DO UPDATE SET detail = excluded.detail," +
			"                  title = excluded.title," +
			"                  location = excluded.location;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(source.ID, source.Detail, source.SourceTitle, source.SourceLocation)
	if err != nil {
		log.Fatal(err)
	}
}

func sqlInsertNote(note note, db *sql.DB) {
	fmt.Printf("Note: id %v Detail %v", note.ID, note.Detail)
	stmt, err := db.Prepare(
		"INSERT INTO note(id, detail)" +
			" VALUES ($1, $2)" +
			" ON CONFLICT (id) " +
			"    DO UPDATE SET detail = excluded.detail;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(note.ID, note.Detail)
	if err != nil {
		log.Fatal(err)
	}
}

func sqlInsertFamily(family family, db *sql.DB) {
	fmt.Printf("Family: id: %v userid %v primeid %v partnerid %v note %v\n",
		family.ID, family.UserID, family.PrimeID.ID, family.PartnerID.ID, family.NoteID.ID)
	stmt, err := db.Prepare(
		"INSERT INTO family(id, userid, primeid, partnerid, noteid)" +
			" VALUES ($1, $2, $3, $4, $5)" +
			" ON CONFLICT (id) " +
			"    DO UPDATE SET userid = excluded.userid," +
			"                  primeid = excluded.primeid," +
			"                  partnerid = excluded.partnerid," +
			"                  noteid = excluded.noteid;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(family.ID, family.UserID,
		NewNullString(family.PrimeID.ID), NewNullString(family.PartnerID.ID),
		NewNullString(family.NoteID.ID))
	if err != nil {
		log.Fatal(err)
	}
}

func sqlInsertPerson(person person, db *sql.DB) {
	fmt.Printf("Person: id: %v userid %v sex %v private %v note %v source %v\n",
		person.ID, person.UserID, person.BirthSex,
		person.IsPrivate, person.NoteID.ID, person.SourceID.ID)
	stmt, err := db.Prepare(
		"INSERT INTO person(id, userid, birthsex, isPrivate, noteid, sourceid)" +
			" VALUES ($1, $2, $3, $4, $5, $6)" +
			" ON CONFLICT (id) " +
			"    DO UPDATE SET userid = excluded.userid," +
			"                  birthsex = excluded.birthsex," +
			"                  isPrivate = excluded.isPrivate," +
			"                  noteid = excluded.noteid," +
			"                  sourceid = excluded.sourceid;")
	if err != nil {
		log.Fatal(err)
	}
	isPrivate := person.IsPrivate == "True"
	birthSex := NewNullString(person.BirthSex)
	noteid := NewNullString(person.NoteID.ID)
	sourceid := NewNullString(person.SourceID.ID)
	_, err = stmt.Exec(person.ID, person.UserID, birthSex, isPrivate, noteid, sourceid)
	if err != nil {
		log.Fatal(err)
	}
}
