package googleplay

import (
   "bytes"
   "fmt"
   "github.com/89z/googleplay"
   "net/url"
   "os"
   "testing"
   "time"
)

func TestUpload(t *testing.T) {
   ac2dmToken, ok := os.LookupEnv("ac2dmToken")
   if ! ok {
      panic("ac2dmToken")
   }
   check, err := newCheckin()
   if err != nil {
      t.Fatal(err)
   }
   deviceID := check.String()
   if err := newDevice().upload(ac2dmToken, deviceID); err != nil {
      t.Fatal(err)
   }
   fmt.Println(deviceID)
   auth := googleplay.Auth{
      url.Values{
         "Auth": {ac2dmToken},
      },
   }
   time.Sleep(16 * time.Second)
   det, err := auth.Details(deviceID, "com.google.android.youtube")
   if err != nil {
      t.Fatal(err)
   }
   vers := []string{"16.", "16.4", "16.43.", "16.43.3", "16.43.34"}
   for _, ver := range vers {
      if bytes.Contains(det, []byte(ver)) {
         fmt.Println("pass", ver)
      } else {
         fmt.Println("fail", ver)
      }
   }
}