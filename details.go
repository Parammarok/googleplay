package googleplay

import (
   "errors"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func (h Header) Details(app string) (*Details, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/details", nil,
   )
   if err != nil {
      return nil, err
   }
   // half of the apps I test require User-Agent,
   // so just set it for all of them
   h.Set_Agent(req.Header)
   h.Set_Auth(req.Header)
   h.Set_Device(req.Header)
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   response_wrapper := make(protobuf.Message)
   response_wrapper.ReadFrom(res.Body)
   // .payload.detailsResponse.docV2
   docV2 := response_wrapper.Get(1).Get(2).Get(4)
   var det Details
   // The following fields will fail with wrong ABI, so try them first. If the
   // first one passes, then use native error for the rest.
   // .details.appDetails.versionCode
   det.Version_Code, err = docV2.Get(13).Get(1).GetVarint(3)
   if err != nil {
      return nil, version_error{app}
   }
   // .details.appDetails.versionString
   det.Version, err = docV2.Get(13).Get(1).GetString(4)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.installationSize
   det.Size, err = docV2.Get(13).Get(1).GetVarint(9)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.uploadDate
   det.Upload_Date, err = docV2.Get(13).Get(1).GetString(16)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.file
   for _, file := range docV2.Get(13).Get(1).GetMessages(17) {
      // .fileType
      typ, err := file.GetVarint(1)
      if err != nil {
         return nil, err
      }
      det.File = append(det.File, typ)
   }
   // The following fields should work with any ABI.
   // .title
   det.Title, err = docV2.GetString(5)
   if err != nil {
      return nil, err
   }
   // .creator
   det.Creator, err = docV2.GetString(6)
   if err != nil {
      return nil, err
   }
   // .offer.micros
   det.Micros, err = docV2.Get(8).GetVarint(1)
   if err != nil {
      return nil, err
   }
   // .offer.currencyCode
   det.Currency_Code, err = docV2.Get(8).GetString(2)
   if err != nil {
      return nil, err
   }
   // I dont know the name of field 70
   // .details.appDetails
   det.Downloads, err = docV2.Get(13).Get(1).GetVarint(70)
   if err != nil {
      return nil, err
   }
   return &det, nil
}

type version_error struct {
   app string
}

func (v version_error) Error() string {
   var buf strings.Builder
   buf.WriteString(v.app)
   buf.WriteString(" versionCode missing\n")
   buf.WriteString("Check nativePlatform")
   return buf.String()
}

type Details struct {
   Creator string
   Currency_Code string
   Downloads uint64
   File []uint64
   Micros uint64
   Size uint64
   Title string
   Upload_Date string // Jun 1, 2021
   Version string
   Version_Code uint64
}

func (d Details) String() string {
   var buf []byte
   buf = append(buf, "Title: "...)
   buf = append(buf, d.Title...)
   buf = append(buf, "\nCreator: "...)
   buf = append(buf, d.Creator...)
   buf = append(buf, "\nUploadDate: "...)
   buf = append(buf, d.Upload_Date...)
   buf = append(buf, "\nVersionString: "...)
   buf = append(buf, d.Version...)
   buf = append(buf, "\nVersionCode: "...)
   buf = strconv.AppendUint(buf, d.Version_Code, 10)
   buf = append(buf, "\nNumDownloads: "...)
   buf = append(buf, format.LabelNumber(d.Downloads)...)
   buf = append(buf, "\nSize: "...)
   buf = append(buf, format.LabelSize(d.Size)...)
   buf = append(buf, "\nFile:"...)
   for _, file := range d.File {
      if file == 0 {
         buf = append(buf, " APK"...)
      } else {
         buf = append(buf, " OBB"...)
      }
   }
   buf = append(buf, "\nOffer: "...)
   buf = strconv.AppendUint(buf, d.Micros, 10)
   buf = append(buf, ' ')
   buf = append(buf, d.Currency_Code...)
   return string(buf)
}

// This only works with English. You can force English with:
// Accept-Language: en
func (d Details) Time() (time.Time, error) {
   return time.Parse("Jan 2, 2006", d.Upload_Date)
}
