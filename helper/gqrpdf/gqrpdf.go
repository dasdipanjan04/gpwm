package gqrpdf

import (
	"github.com/dasdipanjan04/gpwm/helper/glogger"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// Generate a pdf file with QRCoded Master Account Key
func MasterKeyQRCodePDFGenerator(masterKey string, firstName string, lastname string) {

	qrpdf := pdf.NewMaroto(consts.Portrait, consts.Letter)
	qrpdf.Row(12, func() {
		qrpdf.Text(string("Dear Mr./Mrs. "+firstName+" "+lastname),
			props.Text{
				Top:         12,
				Size:        20,
				Extrapolate: true,
				Align:       consts.Left,
				Family:      "arial",
			})
	})
	qrpdf.Row(12, func() {
		qrpdf.Text("Please scan the QR code below to get your masterkey.",
			props.Text{
				Top:         12,
				Size:        20,
				Extrapolate: true,
				Align:       consts.Left,
				Family:      "arial",
			})
	})
	qrpdf.Row(50, func() {
		qrpdf.Text("Please store it in a secure place.",
			props.Text{
				Top:         12,
				Size:        20,
				Extrapolate: true,
				Align:       consts.Left,
				Family:      "arial",
			})
	})
	qrpdf.Row(150, func() {
		qrpdf.Col(12, func() {
			qrpdf.QrCode(masterKey,
				props.Rect{
					Top:     2,
					Center:  true,
					Percent: 400,
				})
		})
	})

	err := qrpdf.OutputFileAndClose("gqrpdf.pdf")
	if err != nil {
		glogger.Glog("helper:gqrpdf:MasterKeyQRCodePDFGenerator", err.Error())
	}

}
