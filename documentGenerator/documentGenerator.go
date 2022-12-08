package documentGenerator

import (
	"fmt"
	"log"
	"strconv"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/tsawlen/matchingAppInvoiceService/common/dataStructures"
)

func GenerateInvoicePDF(billed *dataStructures.User, invoice *dataStructures.Invoice) []byte {
	document := pdf.NewMaroto(consts.Portrait, consts.A4)
	document.SetPageMargins(20, 10, 20)
	document.RegisterHeader(func() {
		document.Row(50, func() {
			document.Col(12, func() {
				document.Text("Finder", props.Text{
					Size:  32,
					Align: consts.Right,
					Style: consts.Bold,
					Top:   3,
					Color: color.NewBlack(),
				})
			})
		})
	})
	buildRecipientSection(document, billed)
	generateInvoiceInformation(document, billed, invoice)
	generateInvoiceTable(document, invoice)
	generateTotalSection(document, invoice)
	registerDocFooter(document)
	if err := document.OutputFileAndClose("test.pdf"); err != nil {
		log.Println(err)
	}
	outDoc, err := document.Output()
	if err != nil {
		log.Println(err)
	}
	return outDoc.Bytes()

}

// Parts builder

func buildRecipientSection(document pdf.Maroto, billed *dataStructures.User) {
	document.Row(20, func() {
		document.Col(7, func() {
			document.Text("Finder Straße 25, 68160 Mannheim", props.Text{
				Size:  6,
				Align: consts.Left,
			})
			document.Text(billed.First_name+" "+billed.Name, props.Text{
				Top:   4,
				Size:  8,
				Align: consts.Left,
			})
			document.Text(billed.Street+" "+billed.HouseNumber, props.Text{
				Top:   7,
				Size:  8,
				Align: consts.Left,
			})
			document.Text(strconv.Itoa(int(billed.City.PLZ))+" "+billed.City.Place, props.Text{
				Top:   10,
				Size:  8,
				Align: consts.Left,
			})
		})
		document.Col(5, func() {
			document.Text("Finder.de", props.Text{
				Size:  10,
				Align: consts.Left,
				Style: consts.Bold,
			})
			document.Text("Straße 25", props.Text{
				Size:  8,
				Align: consts.Left,
				Top:   4,
			})
			document.Text("68160 Mannheim", props.Text{
				Size:  8,
				Align: consts.Left,
				Top:   7,
			})
			document.Text("Tel: XXXX XXXXXXXX", props.Text{
				Size:  8,
				Align: consts.Left,
				Top:   10,
			})
			document.Text("E-Mail: finder.tutoring@gmail.com", props.Text{
				Size:  8,
				Align: consts.Left,
				Top:   13,
			})

		})
	})
}

func generateInvoiceInformation(document pdf.Maroto, billed *dataStructures.User, invoice *dataStructures.Invoice) {
	document.Row(30, func() {
		document.Col(7, func() {
			document.Text("Rechnung", props.Text{
				Top:   10,
				Size:  12,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
		document.Col(5, func() {
			document.Text("Rechnungsnummer: "+invoice.Id.String(), props.Text{
				Top:   7,
				Size:  7,
				Align: consts.Left,
			})
			document.Text("Rechnungsdatum: \t\t\t"+invoice.CreatedAt.Format("02.01.2006"), props.Text{
				Top:   10,
				Size:  7,
				Align: consts.Left,
			})
			document.Text("Kundennummer: \t\t\t\t\t\t"+strconv.Itoa(invoice.Payer), props.Text{
				Top:   13,
				Size:  7,
				Align: consts.Left,
			})
		})
	})
}

func generateInvoiceTable(document pdf.Maroto, invoice *dataStructures.Invoice) {
	tableHeadings := []string{"Menge", "Bezeichnung", "Netto", "Brutto", "Ust."}
	contents := generateContent(invoice)

	document.TableList(tableHeadings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{1, 7, 2, 2, 1},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{1, 7, 2, 2, 1},
		},
		Align:                consts.Left,
		AlternatedBackground: getLightGrayColor(),
		HeaderContentSpace:   1,
		Line:                 false,
	})
}

func generateTotalSection(document pdf.Maroto, invoice *dataStructures.Invoice) {
	document.Row(12, func() {
		document.Col(12, func() {
			document.Text("Summe Netto: \t\t\t\t\t\t"+calcNetto(invoice.Amount)+" EUR", props.Text{
				Top:   4,
				Align: consts.Right,
				Size:  8,
				Style: consts.Bold,
			})
			document.Text("Summe USt.: \t\t\t\t\t\t"+calcTax(invoice.Amount)+" EUR", props.Text{
				Top:   8,
				Align: consts.Right,
				Size:  8,
				Style: consts.Bold,
			})
			document.Text("Endbetrag: \t\t\t\t"+convertToString(invoice.Amount)+" EUR", props.Text{
				Top:   12,
				Align: consts.Right,
				Size:  8,
				Style: consts.Bold,
			})
		})
	})
}

func registerDocFooter(document pdf.Maroto) {
	document.RegisterFooter(func() {
		document.Row(10, func() {
			document.Col(4, func() {
				document.Text("Finder.de", props.Text{
					Size:  6,
					Align: consts.Left,
				})
				document.Text("Straße 25", props.Text{
					Size:  6,
					Align: consts.Left,
					Top:   2,
				})
				document.Text("68160 Mannheim", props.Text{
					Size:  6,
					Align: consts.Left,
					Top:   4,
				})
				document.Text("Tel: XXXX XXXXXXXX", props.Text{
					Size:  6,
					Align: consts.Left,
					Top:   6,
				})
				document.Text("E-Mail: finder.tutoring@gmail.com", props.Text{
					Size:  6,
					Align: consts.Left,
					Top:   8,
				})
			})
			document.Col(4, func() {
				document.Text("Dies ist eine fiktionelle Rechnung\n eines fiktionellen Unternehmens für ein\n Universitäres Programmierprojekt!", props.Text{
					Size:  6,
					Align: consts.Right,
				})
			})
			document.Col(4, func() {
				document.Text("Mock Bank", props.Text{
					Size:  6,
					Align: consts.Right,
				})
				document.Text("IBAN: DE00 0000 0000 0000 00", props.Text{
					Size:  6,
					Align: consts.Right,
					Top:   3,
				})
				document.Text("BIC: MOCKBANK000", props.Text{
					Size:  6,
					Align: consts.Right,
					Top:   6,
				})
			})
		})
	})
}

// Helper

func getLightGrayColor() *color.Color {
	return &color.Color{
		Red:   176,
		Green: 173,
		Blue:  172,
	}
}

func generateContent(invoice *dataStructures.Invoice) [][]string {
	output := [][]string{}
	line := []string{strconv.Itoa(invoice.Hours), invoice.Service, calcNetto(invoice.Amount), convertToString(invoice.Amount), "19%"}
	output = append(output, line)
	return output
}

func calcNetto(brutto float64) string {
	netto := brutto * 0.81
	string := fmt.Sprintf("%.2f", netto)
	return string
}

func calcTax(brutto float64) string {
	netto := brutto * 0.19
	string := fmt.Sprintf("%.2f", netto)
	return string
}

func convertToString(brutto float64) string {
	return fmt.Sprintf("%.2f", brutto)
}
